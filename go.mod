module github.com/ncruces/go-gcp

go 1.16

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10
	go.opencensus.io v0.23.0
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
)

replace github.com/aws/aws-sdk-go => github.com/ncruces/go-gcp/aws-sdk-shim v1.0.0
