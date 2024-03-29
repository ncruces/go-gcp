package glog_test

import "github.com/ncruces/go-gcp/glog"

func init() {
	glog.LogSourceLocation = false
}

func ExamplePrint() {
	glog.Print("Test")
	// Output:
	// {"message":"Test"}
}

func ExampleInfof() {
	glog.Infof("Hello %q!", "Google")
	// Output:
	// {"message":"Hello \"Google\"!","severity":"INFO"}
}

func ExampleWarningj() {
	glog.Warningj("Warning", map[string]string{
		"component": "app",
	})
	// Output:
	// {"component":"app","message":"Warning","severity":"WARNING"}
}

func ExampleWarningw() {
	glog.Warningw("Warning",
		"component", "app")
	// Output:
	// {"component":"app","message":"Warning","severity":"WARNING"}
}
