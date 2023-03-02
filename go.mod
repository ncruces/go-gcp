module github.com/ncruces/go-gcp

go 1.19

require (
	cloud.google.com/go/functions v1.10.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	go.opencensus.io v0.24.0
	golang.org/x/oauth2 v0.5.0
)

require (
	cloud.google.com/go/compute v1.18.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/monitoring v1.12.0 // indirect
	cloud.google.com/go/trace v1.8.0 // indirect
	github.com/aws/aws-sdk-go v1.44.212 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/prometheus/prometheus v0.42.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/api v0.111.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230301171018-9ab4bdc49ad5 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/aws/aws-sdk-go => github.com/ncruces/go-gcp/aws-sdk-shim v1.0.0
