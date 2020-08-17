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

type HTTPFormat struct {
	propagation.HTTPFormat
}
