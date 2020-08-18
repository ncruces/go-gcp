This is a shim to help ensure Google Cloud Trace libraries don't depend on AWS SDK.

To use it, add the following to your `go.mod`:

    replace github.com/aws/aws-sdk-go => github.com/ncruces/go-gcp/aws-sdk-shim v1.0.0