package gcshandler

import (
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
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

	ctx := context.Background()

	if c.Client == nil {
		// use default application credential
		c.Client, _ = storage.NewClient(ctx)
	}

	if c.Client == nil {
		// use anonymous account
		c.Client, _ = storage.NewClient(ctx, option.WithoutAuthentication())
	}

	if c.Client == nil {
		panic("gcshandler: can not init storage client")
	}

	// normalize base path
	c.BasePath = strings.TrimPrefix(c.BasePath, "/")
	c.BasePath = strings.TrimSuffix(c.BasePath, "/")

	bucket := c.Client.Bucket(c.Bucket)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := bucket.Object(strings.TrimPrefix(path.Join(c.BasePath, r.URL.Path), "/"))

		reader, err := obj.NewReader(r.Context())
		if err != nil {
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
