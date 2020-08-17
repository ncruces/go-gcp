package gmutex

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
)

type Mutex struct {
	Bucket  string
	Object  string
	Expires time.Duration
}

var _ sync.Locker = (*Mutex)(nil)

var once sync.Once
var httpClient *http.Client

func initClient(ctx context.Context) (err error) {
	once.Do(func() {
		const scope = "https://www.googleapis.com/auth/devstorage.read_write"
		httpClient, err = google.DefaultClient(ctx, scope)
	})
	return err
}

// Lock locks m. If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m *Mutex) Lock() {
	m.LockContext(context.Background())
}

// Unlock unlocks m. It is a run-time error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() {
	m.UnlockContext(context.Background())
}

func (m *Mutex) LockContext(ctx context.Context) error {
	if err := initClient(ctx); err != nil {
		return err
	}

	url := m.storageURL()
	backoff := 10 * time.Millisecond
	for {
		req, err := http.NewRequestWithContext(ctx, "PUT", url, http.NoBody)
		if err != nil {
			return err
		}
		req.Header.Set("x-goog-if-generation-match", "0")
		res, err := httpClient.Do(req)
		if err == nil {
			res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return nil
			}
		}
		select {
		case <-time.After(backoff):
			backoff *= 2
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *Mutex) UnlockContext(ctx context.Context) error {
	if err := initClient(ctx); err != nil {
		return err
	}
	return nil
}

func (m *Mutex) storageURL() string {
	u := url.URL{
		Scheme: "https",
		Host:   "storage.googleapis.com",
		Path:   m.Bucket + "/" + m.Object,
	}
	return u.String()
}
