package gmutex_test

import (
	"context"
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
	gmutex.HttpClient = http.DefaultClient
	os.Exit(m.Run())
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
