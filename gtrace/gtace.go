// Package gtrace implements tracing for Google Cloud Run and Cloud Functions.
package gtrace

import (
	"runtime"
	"sync"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"

	"github.com/ncruces/go-gcp/glog"
)

var once sync.Once

// Init initializes Cloud Trace.
// Can be called multiple times.
// Logs the error when called asynchronously.
func Init() (err error) {
	_, _, _, caller := runtime.Caller(2)

	once.Do(func() {
		exporter, ierr := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: glog.ProjectID,
		})
		if ierr == nil {
			trace.RegisterExporter(exporter)
		} else if !caller {
			glog.Critical(ierr)
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
