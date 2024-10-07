package gmutex_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ncruces/go-gcp/gmutex"
)

var bucket = os.Getenv("BUCKET")
var object = os.Getenv("OBJECT")

func TestMain(m *testing.M) {
	gmutex.HTTPClient = http.DefaultClient
	if os.Getenv("STORAGE_HOST_EMULATOR") != "" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if bucket != "" && object != "" {
		os.Exit(m.Run())
	}
}

func TestMutex_contention(t *testing.T) {
	ctx := context.Background()

	var failed bool
	var running bool
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mtx, err := gmutex.New(ctx, bucket, object, 5*time.Minute)
			if err != nil {
				t.Error(err)
				return
			}

			t.Log("locking", i)
			if err := mtx.Lock(ctx); err != nil {
				t.Error(err)
				return
			}
			t.Log("locked", i)

			if running || failed {
				failed = true
			} else {
				running = true
				time.Sleep(time.Second)
				running = false
			}

			t.Log("unlocking", i)
			if err := mtx.Unlock(ctx); err != nil {
				t.Error(err)
				return
			}
			t.Log("unlocked", i)
		}(i)
	}
	wg.Wait()

	if failed {
		t.Fail()
	}
}

func TestMutex_expiration(t *testing.T) {
	ctx := context.Background()
	mtx, err := gmutex.New(ctx, bucket, object, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("locking")
	if err := mtx.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("locked")
	mtx.Abandon()
	t.Log("abandoned")

	t.Log("locking")
	if err := mtx.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("locked")

	t.Log("unlocking")
	if err := mtx.Unlock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("unlocked")
}

func TestMutex_extension(t *testing.T) {
	ctx := context.Background()
	mtx, err := gmutex.New(ctx, bucket, object, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("locking")
	if err := mtx.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("locked")

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)

		t.Log("extending")
		if err := mtx.Extend(ctx); err != nil {
			t.Fatal(err)
		}
		t.Log("extended")
	}

	mtx.Abandon()
	t.Log("abandoned")

	t.Log("locking")
	if err := mtx.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("locked")

	t.Log("unlocking")
	if err := mtx.Unlock(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("unlocked")
}

func TestMutex_SetTTL(t *testing.T) {
	tests := []struct {
		name string
		ttl  time.Duration
		want time.Duration
	}{
		{"zero", 0, 0},
		{"zero", -1, 0},
		{"zero", +1, time.Second},
		{"nano", time.Nanosecond, time.Second},
		{"micro", time.Microsecond, time.Second},
		{"milli", time.Millisecond, time.Second},
		{"one", time.Second, time.Second},
		{"one", time.Second - 1, time.Second},
		{"one", time.Second + 1, 2 * time.Second},
		{"negative", -time.Hour, 0},
	}

	ctx := context.Background()
	mtx, err := gmutex.New(ctx, bucket, object, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mtx.SetTTL(tt.ttl)
			got := mtx.TTL()
			if got != tt.want {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
