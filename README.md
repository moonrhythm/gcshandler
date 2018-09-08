# gcshandler

[![Build Status](https://travis-ci.com/moonrhythm/gcshandler.svg?branch=master)](https://travis-ci.com/moonrhythm/gcshandler)
[![Coverage Status](https://coveralls.io/repos/github/moonrhythm/gcshandler/badge.svg?branch=master)](https://coveralls.io/github/moonrhythm/gcshandler?branch=master)
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
