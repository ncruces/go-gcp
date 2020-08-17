module github.com/ncruces/go-gcp

go 1.11

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.3
	go.opencensus.io v0.22.4
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

replace github.com/aws/aws-sdk-go => github.com/ncruces/go-gcp/aws-sdk-shim v1.0.0
