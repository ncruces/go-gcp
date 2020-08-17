# A global mutex using [Google Cloud Storage](https://cloud.google.com/storage)

[![PkgGoDev](https://pkg.go.dev/badge/image)](https://pkg.go.dev/github.com/ncruces/go-gcp/gmutex)

Work in Progress, based on [github.com/marcacohen/gcslock](https://github.com/marcacohen/gcslock).

Major goal is to have locks that expire automatically after some time.
Other improvements include reusing the HTTP client between locks.