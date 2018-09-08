# gcshandler

[![Build Status](https://travis-ci.org/moonrhythm/gcshandler.svg?branch=master)](https://travis-ci.org/moonrhythm/gcshandler)
[![codecov](https://codecov.io/gh/moonrhythm/gcshandler/branch/master/graph/badge.svg)](https://codecov.io/gh/moonrhythm/gcshandler)
[![Go Report Card](https://goreportcard.com/badge/github.com/moonrhythm/gcshandler)](https://goreportcard.com/report/github.com/moonrhythm/gcshandler)
[![GoDoc](https://godoc.org/github.com/moonrhythm/gcshandler?status.svg)](https://godoc.org/github.com/moonrhythm/gcshandler)

## Example

```go
m.Handle("/-/", http.StripPrefix("/-", gcshandler.New(gcshandler.Config{
    Client: nil, // *storage.Client
    Bucket: "acoshift-test",
    BasePath: "folder",
    CacheControl: "public, max-age=7200",
    Fallback: webstatic.New("assets"), // github.com/acoshift/webstatic
})))
```

## License

MIT
