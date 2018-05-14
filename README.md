# gcshandler

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
