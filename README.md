# gcshandler

[![Build Status](https://travis-ci.org/acoshift/gcshandler.svg?branch=master)](https://travis-ci.org/acoshift/gcshandler)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/gcshandler/badge.svg?branch=master)](https://coveralls.io/github/acoshift/gcshandler?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/gcshandler)](https://goreportcard.com/report/github.com/acoshift/gcshandler)
[![GoDoc](https://godoc.org/github.com/acoshift/gcshandler?status.svg)](https://godoc.org/github.com/acoshift/gcshandler)

## Example

```go
m.Handle("/-/", http.StripPrefix("/-", cacheControl(gcshandler.New(gcshandler.Config{
    Bucket:   "acoshift-test",
    BasePath: "/folder",
    Fallback: webstatic.New("assets"), // github.com/acoshift/webstatic
    ModifyResponse: func(w *http.Response) error {
        w.Header.Del("Cache-Control")
        return nil
    },
}))))
```
