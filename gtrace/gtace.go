// Package gtrace implements tracing for Google Cloud Run and Cloud Functions.
package gtrace

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"sync"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var once sync.Once

// ProjectID should be set to the Google Cloud project ID.
var ProjectID string = os.Getenv("GOOGLE_CLOUD_PROJECT")

// Init initializes Cloud Trace.
// Can be called multiple times.
// Logs the error if called asynchronously.
func Init() (err error) {
	callers := runtime.Callers(3, make([]uintptr, 1))

	once.Do(func() {
		exporter, ierr := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: ProjectID,
		})
		if ierr == nil {
			trace.RegisterExporter(exporter)
			return
		}
		if callers == 0 {
			json.NewEncoder(os.Stderr).Encode(map[string]string{
				"message":  ierr.Error(),
				"severity": "CRITICAL",
			})
		}
		err = ierr
	})

	return
}

// HTTPFormat implements propagation.HTTPFormat to propagate traces in
// HTTP headers for Cloud Trace.
type HTTPFormat struct {
	propagation.HTTPFormat
}

// NewHTTPClient returns a tracing http.Client.
func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
		},
	}
}

// NewHTTPTransport returns a tracing http.RoundTripper.
func NewHTTPTransport() http.RoundTripper {
	return &ochttp.Transport{
		// Use Google Cloud propagation format.
		Propagation: &propagation.HTTPFormat{},
	}
}

// NewHTTPHandler returns a tracing http.Handler.
func NewHTTPHandler() http.Handler {
	return &ochttp.Handler{
		// Use the Google Cloud propagation format.
		Propagation: &propagation.HTTPFormat{},
	}
}
