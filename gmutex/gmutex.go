package gmutex

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	ttl time.Duration
	gen string // mutable state
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
	return &Mutex{
		url: url.String(),
		ttl: ttl,
	}, nil
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
	var backoff expBackOff
	for i := 0; ; i++ {
		// Create the lock object, if it does not yet exist.
		status, _ := m.createObject(ctx, "0", data)
		if status == http.StatusOK {
			// Lock acquired.
			return nil
		}
		// If the lock object already existed, check if it is expired.
		if status == http.StatusPreconditionFailed && m.expireObject(ctx) {
			// Lock expired, retry immediately.
			continue
		}
		// If the lock is contended, wait for it to expire.
		// This is much cheaper than repeatedly trying to acquire it.
		if status == http.StatusTooManyRequests || i > 2 {
			i = 0
			for {
				// Exponential backoff.
				if err := backoff.wait(ctx); err != nil {
					return err
				}
				if m.expireObject(ctx) {
					break
				}
			}
			// Lock expired, retry immediately.
			continue
		}
		// Exponential backoff.
		if err := backoff.wait(ctx); err != nil {
			return err
		}
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
	var backoff expBackOff
	for i := 0; ; i++ {
		// Create the lock object, if it does not yet exist.
		status, err := m.createObject(ctx, "0", data)
		if status == http.StatusOK {
			// Lock acquired.
			return true, nil
		}
		// If the lock object already existed, check if it is expired.
		if status == http.StatusPreconditionFailed {
			if m.expireObject(ctx) {
				// Lock expired, retry immediately.
				continue
			}
			// Lock held, give up.
			return false, nil
		}
		// If the lock is contended, give up.
		if status == http.StatusTooManyRequests || i > 2 {
			return false, err
		}
		// Exponential backoff.
		if err := backoff.wait(ctx); err != nil {
			return false, err
		}
	}
}

// Unlock unlocks m, deleting any attached data.
// Returns an error if the lock had already expired, and mutual
// exclusion was not ensured.
func (m *Mutex) Unlock(ctx context.Context) error {
	if m.gen == "" {
		panic("gmutex: unlock of unlocked mutex")
	}
	var backoff linBackOff
	for {
		// Delete the lock object, if we still own it.
		status := m.deleteObject(ctx, m.gen)
		if status == http.StatusOK || status == http.StatusNoContent {
			return nil
		}
		// If we no longer owned it, or it doesn't exist, report error.
		if status == http.StatusPreconditionFailed || status == http.StatusNotFound {
			return errors.New("failed to unlock mutex: stale lock")
		}
		// Linear backoff.
		if err := backoff.wait(ctx); err != nil {
			return err
		}
	}
}

// Update updates attached data, extending the expiration time of m.
// Returns an error if the lock has already expired, and mutual
// exclusion can not be ensured.
func (m *Mutex) Update(ctx context.Context, data io.Reader) error {
	if m.gen == "" {
		panic("gmutex: update of unlocked mutex")
	}
	var backoff linBackOff
	for {
		// Update the lock object, if we still own it.
		status, _ := m.createObject(ctx, m.gen, data)
		if status == http.StatusOK {
			return nil
		}
		// If we no longer owned it, or it doesn't exist, abort.
		if status == http.StatusPreconditionFailed || status == http.StatusNotFound {
			return errors.New("failed to extend mutex: stale lock, abort")
		}
		// Linear backoff.
		if err := backoff.wait(ctx); err != nil {
			return err
		}
	}
}

// Inspect inspects m, returning its locked state and fetching attached data.
func (m *Mutex) Inspect(ctx context.Context, data io.Writer) (bool, error) {
	var backoff expBackOff
	for i := 0; ; i++ {
		// Inspect the lock object.
		status, err := m.inspectObject(ctx, data)
		if status == http.StatusOK {
			return true, err
		}
		if status == http.StatusNotFound {
			return false, err
		}
		// Exponential backoff.
		if err := backoff.wait(ctx); err != nil {
			return false, err
		}
	}
}

func (m *Mutex) createObject(ctx context.Context, generation string, data io.Reader) (int, error) {
	// Create/update the lock object if the generation matches.
	req, err := http.NewRequestWithContext(ctx, "PUT", m.url, data)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("x-goog-if-generation-match", generation)
	req.Header.Set("x-goog-meta-ttl", strconv.FormatInt(int64(m.ttl/time.Second), 10))

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	if res.StatusCode == http.StatusOK {
		m.gen = res.Header.Get("x-goog-generation")
	}
	return res.StatusCode, nil
}

func (m *Mutex) deleteObject(ctx context.Context, generation string) int {
	// Delete the lock object if the generation matches.
	req, err := http.NewRequestWithContext(ctx, "DELETE", m.url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("x-goog-if-generation-match", generation)

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0
	}
	res.Body.Close()
	return res.StatusCode
}

func (m *Mutex) inspectObject(ctx context.Context, data io.Writer) (int, error) {
	var method string
	if data == nil {
		method = "HEAD"
	}

	// Inspect the lock object's status.
	req, err := http.NewRequestWithContext(ctx, method, m.url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-cache")

	res, err := HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// If it exists, but is expired, act as if it didn't.
	if res.StatusCode == http.StatusOK && checkExpired(res) {
		return http.StatusNotFound, nil
	}

	if res.StatusCode == http.StatusOK && data != nil {
		_, err := io.Copy(data, res.Body)
		return res.StatusCode, err
	}
	return res.StatusCode, nil
}

func (m *Mutex) expireObject(ctx context.Context) bool {
	// Inspect the lock object's status.
	req, err := http.NewRequestWithContext(ctx, "HEAD", m.url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cache-Control", "no-cache")

	res, err := HttpClient.Do(req)
	if err != nil {
		return false
	}
	res.Body.Close()
	// If it no longer exists, it's as good as expired.
	if res.StatusCode == http.StatusNotFound {
		return true
	}
	// If it still exists, but is expired, delete it.
	if res.StatusCode == http.StatusOK && checkExpired(res) {
		status := m.deleteObject(ctx, res.Header.Get("x-goog-generation"))
		return status == http.StatusOK || status == http.StatusNoContent || status == http.StatusNotFound
	}
	return false
}

func checkExpired(res *http.Response) bool {
	// Check for expiration using server date.
	ttl, err := strconv.ParseInt(res.Header.Get("x-goog-meta-ttl"), 10, 64)
	if err != nil || ttl <= 0 {
		return false
	}
	now, err := http.ParseTime(res.Header.Get("Date"))
	if err != nil {
		return false
	}
	modifed, err := http.ParseTime(res.Header.Get("Last-Modified"))
	if err != nil {
		return false
	}
	expires := modifed.Add(time.Duration(ttl) * time.Second)
	return expires.Before(now)
}
