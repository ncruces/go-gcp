package glog_test

import (
	"log"

	"github.com/ncruces/go-gcp/glog"
)

func ExampleSetupLogger() {
	glog.SetupLogger(log.Default())
	log.Print("Test")
	// Output:
	// {"message":"Test"}
}
