package gcshandler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Config is gcshandler config
type Config struct {
	Bucket   string
	BasePath string
	Fallback http.Handler
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

		return nil
	}

	rev := &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyResponse,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := &bufferResponseWriter{}
		rev.ServeHTTP(nw, r)

		if nw.statusCode >= 400 {
			c.Fallback.ServeHTTP(w, r)
			return
		}

		h := nw.Header()
		hh := w.Header()
		for k, v := range h {
			for _, vv := range v {
				hh.Add(k, vv)
			}
		}
		w.WriteHeader(nw.statusCode)
		io.Copy(w, &nw.buf)
	})
}

type bufferResponseWriter struct {
	buf         bytes.Buffer
	wroteHeader bool
	statusCode  int
	header      http.Header
}

func (w *bufferResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.statusCode = code
}

func (w *bufferResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *bufferResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.buf.Write(p)
}
