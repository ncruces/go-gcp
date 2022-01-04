// Package gmutex implements a global mutex using Google Cloud Storage.
package gmutex

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/oauth2/google"
)

// HTTPClient should be set to an http.Client before first use.
// If unset google.DefaultClient will be used.
var HTTPClient *http.Client

var initMtx sync.Mutex

func initClient(ctx context.Context) (err error) {
	initMtx.Lock()
	defer initMtx.Unlock()
	if HTTPClient == nil {
		const scope = "https://www.googleapis.com/auth/devstorage.read_write"
		HTTPClient, err = google.DefaultClient(ctx, scope)
	}
	return err
}
