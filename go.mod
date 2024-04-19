module github.com/ncruces/go-gcp

go 1.19

require (
	cloud.google.com/go/functions v1.16.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	go.opencensus.io v0.24.0
	golang.org/x/oauth2 v0.17.0
)

require (
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/monitoring v1.18.0 // indirect
	cloud.google.com/go/trace v1.10.5 // indirect
	github.com/aws/aws-sdk-go v1.50.26 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.2 // indirect
	github.com/prometheus/prometheus v0.50.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.167.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240221002015-b0ce06bbee7c // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/aws/aws-sdk-go => github.com/ncruces/go-gcp/aws-sdk-shim v1.0.0
