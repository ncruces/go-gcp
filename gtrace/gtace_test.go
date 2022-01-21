package gtrace_test

import (
	"net/http"
	"os"

	"github.com/ncruces/go-gcp/glog"
	"github.com/ncruces/go-gcp/gtrace"
)

func Example() {
	go gtrace.Init()
	glog.Notice("Starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", http.NotFound)

	glog.Critical(http.ListenAndServe(":"+port, gtrace.NewHTTPHandler()))
}
