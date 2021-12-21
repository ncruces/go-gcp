// Package gmutex implements a global mutex using Google Cloud Storage.
package gmutex

import (
	"context"
	"net/http"
	"sync"

	"golang.org/x/oauth2/google"
)

// HttpClient should be set to an http.Client before first use.
// If unset google.DefaultClient will be used.
var HttpClient *http.Client

var initMtx sync.Mutex

func initClient(ctx context.Context) (err error) {
	initMtx.Lock()
	defer initMtx.Unlock()
	if HttpClient == nil {
		const scope = "https://www.googleapis.com/auth/devstorage.read_write"
		HttpClient, err = google.DefaultClient(ctx, scope)
	}
	return err
}
