package gmutex

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// A Mutex is a global, mutual exclusion lock
// that uses an object in Google Cloud Storage
// to serialize computations across the internet.
//
// A Mutex can optionally have data attached to it while it is held.
// While there is no limit to the size of this data,
// it is best kept small.
// Provided data must be of type *bytes.Buffer, *bytes.Reader,
// or *strings.Reader.
//
// Given the latency and scalability properties of Google Cloud Storage,
// a Mutex is best used to serialize long-running, high-latency
// compute processes.
// Critical sections should span seconds.
// Expect an uncontended mutex to take tens to hundreds of milliseconds
// to acquire, and a contended one multiple seconds after release.
//
// An instance of Mutex is not associated with a particular goroutine
// (it is allowed for one goroutine to lock a Mutex
// and then arrange for another goroutine to unlock it),
// but it is not safe for concurrent use by multiple goroutines.
type Mutex struct {
	_   noCopy
	url string
	ttl int64
	gen string
}

// New creates a new Mutex at the given bucket and object,
// with the given time-to-live.
func New(ctx context.Context, bucket, object string, ttl time.Duration) (*Mutex, error) {
	if err := initClient(ctx); err != nil {
		return nil, err
	}
	url := url.URL{
		Scheme: "https",
		Host:   "storage.googleapis.com",
		Path:   bucket + "/" + object,
	}
	var m Mutex
	m.SetTTL(ttl)
	m.url = url.String()
	return &m, nil
}

// TTL gets the time-to-live to use when the mutex is locked or updated.
func (m *Mutex) TTL() time.Duration {
	return time.Duration(m.ttl) * time.Second
}

// SetTTL sets the time-to-live to use when the mutex is locked or updated.
// The time-to-live is rounded up to the nearest second.
// Negative or zero time-to-live means the lock never expires.
func (m *Mutex) SetTTL(ttl time.Duration) {
	ttl += time.Second - time.Nanosecond
	if ttl > 0 {
		m.ttl = int64(ttl / time.Second)
	} else {
		m.ttl = 0
	}
}

// Locker gets a Locker that uses context.Background to call Lock and Unlock,
// and panics on error.
func (m *Mutex) Locker() sync.Locker {
	return locker{m}
}

// Lock locks m.
// If the lock is already in use,
// the calling goroutine blocks until the mutex is available,
// or the context expires.
// Returns nil if the lock was taken successfully.
func (m *Mutex) Lock(ctx context.Context) error {
	return m.LockData(ctx, nil)
}

// LockData locks m with attached data.
// If the lock is already in use,
// the calling goroutine blocks until the mutex is available,
// or the context expires.
// Returns nil if the lock was taken successfully
// (and the attached data stored).
func (m *Mutex) LockData(ctx context.Context, data io.Reader) error {
	if m.gen != "" {
		panic("gmutex: lock of locked mutex")
	}
	if !rewindable(data) {
		panic("gmutex: data not rewindable")
	}

	var backoff expBackOff // Exponential backoff because we don't hold the lock.
	generation := ""       // Empty generation because we expect the lock not to exist.

	for {
		// Create the lock object, at the expected generation.
		status, gen, err := m.createObject(ctx, generation, data)
		if status == http.StatusOK {
			// Lock acquired.
			m.gen = gen
			return nil
		}

		// If the lock object exists at another generation, let's inspect it.
		if status == http.StatusPreconditionFailed {
			status, gen, err = m.inspectObject(ctx, nil)
		}
		// While the lock object exists, and for transient errors, backoff and retry.
		for status == http.StatusOK || retriable(status, err) {
			if err := backoff.wait(ctx); err != nil {
				return err
			}
			status, gen, err = m.inspectObject(ctx, nil)
		}
		// If the lock object no longer exists, or has expired, we can acquire it.
		if status == http.StatusNotFound {
			generation = gen
			continue
		}

		// Can't recover, give up.
		if err != nil {
			return fmt.Errorf("lock mutex: %w", err)
		}
		return fmt.Errorf("lock mutex: http status %d: %s", status, http.StatusText(status))
	}
}

// TryLock tries to lock m.
// Returns true if the lock was taken successfully,
// false if the lock is already in use.
func (m *Mutex) TryLock(ctx context.Context) (bool, error) {
	return m.TryLockData(ctx, nil)
}

// TryLockData tries to lock m with attached data.
// Returns true if the lock was taken successfully
// (and the attached data stored),
// false if the lock is already in use.
func (m *Mutex) TryLockData(ctx context.Context, data io.Reader) (bool, error) {
	if m.gen != "" {
		panic("gmutex: lock of locked mutex")
	}
	if !rewindable(data) {
		panic("gmutex: data not rewindable")
	}

	var backoff expBackOff // Exponential backoff because we don't hold the lock.
	generation := ""       // Empty generation because we expect the lock not to exist.

	for {
		// Create the lock object, at the expected generation.
		status, gen, err := m.createObject(ctx, generation, data)
		if status == http.StatusOK {
			// Lock acquired.
			m.gen = gen
			return true, nil
		}

		// If the lock object exists at another generation, let's inspect it.
		if status == http.StatusPreconditionFailed {
			status, gen, err = m.inspectObject(ctx, nil)
		}
		// For transient errors, backoff and retry.
		for retriable(status, err) {
			if err := backoff.wait(ctx); err != nil {
				return false, err
			}
			status, gen, err = m.inspectObject(ctx, nil)
		}
		// If the lock object no longer exists, or has expired, we can acquire it.
		if status == http.StatusNotFound {
			generation = gen
			continue
		}
		// If the lock object exists.
		if status == http.StatusOK {
			// Lock held, give up.
			return false, nil
		}

		// Can't recover, give up.
		if err != nil {
			return false, fmt.Errorf("lock mutex: %w", err)
		}
		return false, fmt.Errorf("lock mutex: http status %d: %s", status, http.StatusText(status))
	}
}

// Unlock unlocks m, deleting any attached data.
// Returns an error if the lock had already expired,
// and mutual exclusion was not ensured.
func (m *Mutex) Unlock(ctx context.Context) error {
	if m.gen == "" {
		panic("gmutex: unlock of unlocked mutex")
	}

	var backoff linBackOff // Linear backoff because we hold the lock.

	for {
		// Delete the lock object, at the expected generation.
		status, err := m.deleteObject(ctx, m.gen)
		if status == http.StatusOK || status == http.StatusNoContent {
			m.gen = ""
			return nil
		}

		// If the lock object exists at another generation, or no longer exists, it is stale.
		if status == http.StatusPreconditionFailed || status == http.StatusNotFound {
			return errors.New("unlock mutex: stale lock")
		}

		// For transient errors, backoff and retry.
		if retriable(status, err) {
			if err := backoff.wait(ctx); err != nil {
				return err
			}
			continue
		}

		// Can't recover, give up.
		if err != nil {
			return fmt.Errorf("unlock mutex: %w", err)
		}
		return fmt.Errorf("unlock mutex: http status %d: %s", status, http.StatusText(status))
	}
}

// Update updates attached data, extending the expiration time of m.
// Returns an error if the lock has already expired,
// and mutual exclusion can not be ensured.
func (m *Mutex) Update(ctx context.Context, data io.Reader) error {
	if m.gen == "" {
		panic("gmutex: update of unlocked mutex")
	}
	if !rewindable(data) {
		panic("gmutex: data not rewindable")
	}

	var backoff linBackOff // Linear backoff because we hold the lock.

	for {
		// Update the lock object, at the expected generation.
		status, gen, err := m.createObject(ctx, m.gen, data)
		if status == http.StatusOK {
			// Lock updated.
			m.gen = gen
			return nil
		}

		// If the lock object exists at another generation, or no longer exists, it is stale.
		if status == http.StatusPreconditionFailed || status == http.StatusNotFound {
			return errors.New("update mutex: stale lock, abort")
		}

		// For transient errors, backoff and retry.
		if retriable(status, err) {
			if err := backoff.wait(ctx); err != nil {
				return err
			}
			continue
		}

		// Can't recover, give up.
		if err != nil {
			return fmt.Errorf("update mutex: %w", err)
		}
		return fmt.Errorf("update mutex: http status %d: %s", status, http.StatusText(status))
	}
}

// Inspect inspects m, returning its locked state and fetching attached data.
func (m *Mutex) Inspect(ctx context.Context, data io.Writer) (bool, error) {
	var backoff expBackOff // Exponential backoff because we don't hold the lock.

	for {
		// Inspect the lock object.
		status, _, err := m.inspectObject(ctx, data)
		if status == http.StatusOK {
			return true, nil
		}
		if status == http.StatusNotFound {
			return false, nil
		}

		// For transient errors, backoff and retry.
		if retriable(status, err) {
			if err := backoff.wait(ctx); err != nil {
				return false, err
			}
			continue
		}

		// Can't recover, give up.
		if err != nil {
			return false, fmt.Errorf("inspect mutex: %w", err)
		}
		return false, fmt.Errorf("inspect mutex: http status %d: %s", status, http.StatusText(status))
	}
}

// Abandon abandons m, returning a lock id that can be used to call Adopt.
func (m *Mutex) Abandon() string {
	if m.gen == "" {
		panic("gmutex: abandon of unlocked mutex")
	}

	gen := m.gen
	m.gen = ""
	return gen
}

// Adopt adopts an abandoned lock into m,
// and calls Update to ensure mutual exclusion.
func (m *Mutex) Adopt(ctx context.Context, id string, data io.Reader) error {
	if m.gen != "" {
		panic("gmutex: adopt on locked mutex")
	}
	if id == "" || id == "0" {
		panic("gmutex: adopt of invalid lock")
	}

	m.gen = id
	return m.Update(ctx, data)
}

func (m *Mutex) createObject(ctx context.Context, generation string, data io.Reader) (int, string, error) {
	if generation == "" {
		generation = "0"
	}

	// Create/update the lock object if the generation matches.
	req, err := http.NewRequestWithContext(ctx, "PUT", m.url, data)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("x-goog-if-generation-match", generation)
	req.Header.Set("x-goog-meta-ttl", strconv.FormatInt(m.ttl, 10))

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	res.Body.Close()
	return res.StatusCode, res.Header.Get("x-goog-generation"), nil
}

func (m *Mutex) deleteObject(ctx context.Context, generation string) (int, error) {
	// Delete the lock object if the generation matches.
	req, err := http.NewRequestWithContext(ctx, "DELETE", m.url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("x-goog-if-generation-match", generation)

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

func (m *Mutex) inspectObject(ctx context.Context, data io.Writer) (int, string, error) {
	var method string
	if data == nil {
		method = "HEAD"
	}

	// Get the lock object's status.
	req, err := http.NewRequestWithContext(ctx, method, m.url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-cache")

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	// If it exists, but is expired, act as if it didn't.
	if res.StatusCode == http.StatusOK && expired(res) {
		res.StatusCode = http.StatusNotFound
	}
	if res.StatusCode == http.StatusOK && data != nil {
		_, err = io.Copy(data, res.Body)
	}
	return res.StatusCode, res.Header.Get("x-goog-generation"), nil
}

func retriable(status int, err error) bool {
	// Retry on temporary errors and timeouts.
	if err != nil {
		uerr := url.Error{Err: err}
		return uerr.Temporary() || uerr.Timeout()
	}
	return status == http.StatusTooManyRequests ||
		status == http.StatusRequestTimeout ||
		status == http.StatusInternalServerError ||
		status == http.StatusServiceUnavailable ||
		status == http.StatusBadGateway ||
		status == http.StatusGatewayTimeout
}

func rewindable(body io.Reader) bool {
	switch body.(type) {
	case nil, *bytes.Buffer, *bytes.Reader, *strings.Reader:
		return true
	default:
		return body == http.NoBody
	}
}

func expired(res *http.Response) bool {
	// Check for expiration using server date.
	now, err := http.ParseTime(res.Header.Get("Date"))
	if err != nil {
		return false
	}
	modified, err := http.ParseTime(res.Header.Get("Last-Modified"))
	if err != nil {
		return false
	}
	expiration, err := http.ParseTime(res.Header.Get("x-goog-expiration"))
	if err != nil || true {
		expiration = now
	}
	ttl, err := strconv.ParseInt(res.Header.Get("x-goog-meta-ttl"), 10, 64)
	if err != nil || ttl <= 0 {
		return false
	}
	expires := modified.Add(time.Duration(ttl) * time.Second)
	return expires.Before(now) || expiration.Before(now)
}
