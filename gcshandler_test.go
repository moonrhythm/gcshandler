package gcshandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonrhythm/gcshandler"
	"github.com/stretchr/testify/assert"
)

const bucket = "acoshift-test"

type fallbackHandler struct {
	called bool
}

func (h *fallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	w.Write([]byte("fallback"))
}

func TestHandlerSuccess(t *testing.T) {
	t.Run("RootBasePath", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "/",
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.False(t, fallback.called)
	})

	t.Run("RootBaseNestedPath", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "/",
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/folder/file1", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.False(t, fallback.called)
	})

	t.Run("EmptyBasePath", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "",
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.False(t, fallback.called)
	})

	t.Run("NestedBasePath", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "/folder",
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/file1", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.False(t, fallback.called)
	})

	t.Run("CacheControl", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:       bucket,
			BasePath:     "/",
			Fallback:     &fallback,
			CacheControl: "public, max-age=123",
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "public, max-age=123", w.Header().Get("Cache-Control"))
	})
}

func TestHandlerNotFound(t *testing.T) {
	t.Run("WithFallback", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "/",
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/not-exists-file", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.True(t, fallback.called)
		assert.Equal(t, "fallback", w.Body.String())
	})

	t.Run("WithoutFallback", func(t *testing.T) {
		h := gcshandler.New(gcshandler.Config{
			Bucket:   bucket,
			BasePath: "/",
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/not-exists-file", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 404, w.Code)
	})
}

func TestEmptyBucket(t *testing.T) {
	t.Run("WithFallback", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Fallback: &fallback,
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/not-exists-file", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 200, w.Code)
		assert.True(t, fallback.called)
		assert.Equal(t, "fallback", w.Body.String())
	})

	t.Run("WithoutFallback", func(t *testing.T) {
		h := gcshandler.New(gcshandler.Config{})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/not-exists-file", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 404, w.Code)
	})
}
