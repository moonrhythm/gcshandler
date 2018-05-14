package gcshandler

import (
	"net/http"
	"net/http/httputil"
	"strings"
)

// Config is gcshandler config
type Config struct {
	Bucket         string
	BasePath       string
	Fallback       http.Handler
	ModifyResponse func(*http.Response) error
}

const gcsHost = "storage.googleapis.com"

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

	// default ModifyResponse
	if c.ModifyResponse == nil {
		c.ModifyResponse = func(*http.Response) error { return nil }
	}

	// normalize base path
	if !strings.HasPrefix(c.BasePath, "/") {
		c.BasePath = "/" + c.BasePath
	}
	c.BasePath = strings.TrimSuffix(c.BasePath, "/")
	c.BasePath = "/" + c.Bucket + c.BasePath

	// setup reverse proxy
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	director := func(r *http.Request) {
		r.Host = gcsHost
		r.URL.Scheme = "https"
		r.URL.Host = gcsHost
		r.URL.Path = c.BasePath + "/" + strings.TrimPrefix(r.URL.Path, "/")

		// prevent default user-agent
		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "")
		}

		// remove headers
		r.Header.Del("Cookie")
		r.Header.Del("Accept-Encoding")
	}

	modifyResponse := func(w *http.Response) error {
		w.Header.Del("x-goog-generation")
		w.Header.Del("x-goog-metageneration")
		w.Header.Del("x-goog-stored-content-encoding")
		w.Header.Del("x-goog-stored-content-length")
		w.Header.Del("x-goog-hash")
		w.Header.Del("x-goog-storage-class")
		w.Header.Del("x-goog-meta-goog-reserved-file-mtime")
		w.Header.Del("x-guploader-uploadid")
		w.Header.Del("Alt-Svc")
		w.Header.Del("Server")
		w.Header.Del("Age")

		return c.ModifyResponse(w)
	}

	rev := &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyResponse,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := &responseWriter{
			ResponseWriter: w,
		}
		rev.ServeHTTP(nw, r)

		if nw.fallback {
			c.Fallback.ServeHTTP(w, r)
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	fallback    bool
	header      http.Header
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	if code >= 400 {
		w.fallback = true
		return
	}

	h := w.ResponseWriter.Header()
	for k, v := range w.header {
		for _, vv := range v {
			h.Add(k, vv)
		}
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.fallback {
		return len(p), nil
	}
	return w.ResponseWriter.Write(p)
}
