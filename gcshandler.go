package gcshandler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
)

// Config is gcshandler config
type Config struct {
	Client       *storage.Client
	CacheControl string
	Bucket       string
	BasePath     string
	Fallback     http.Handler
}

// New creates new gcshandler
func New(c Config) http.Handler {
	// default fallback
	if c.Fallback == nil {
		c.Fallback = http.NotFoundHandler()
	}

	// short-circit no bucket
	if c.Bucket == "" {
		return c.Fallback
	}

	if c.Client == nil {
		c.Client, _ = storage.NewClient(context.Background())
	}

	// normalize base path
	c.BasePath = strings.TrimPrefix(c.BasePath, "/")
	c.BasePath = strings.TrimSuffix(c.BasePath, "/")

	bucket := c.Client.Bucket(c.Bucket)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := bucket.Object(strings.TrimPrefix(path.Join(c.BasePath, r.URL.Path), "/"))

		reader, err := obj.NewReader(r.Context())
		if err != nil {
			fmt.Println(err)
			c.Fallback.ServeHTTP(w, r)
			return
		}
		defer reader.Close()

		h := w.Header()
		if v := reader.ContentType(); v != "" {
			h.Set("Content-Type", v)
		}

		if c.CacheControl != "" {
			h.Set("Cache-Control", c.CacheControl)
		} else if v := reader.CacheControl(); v != "" {
			h.Set("Cache-Control", v)
		}

		io.Copy(w, reader)
	})
}
