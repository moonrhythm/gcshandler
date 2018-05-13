package gcshandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/gcshandler"
	"github.com/stretchr/testify/assert"
)

type fallbackHandler struct {
	called bool
}

func (h *fallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	w.Write([]byte("fallback"))
}

func TestHandlerSuccess(t *testing.T) {
	fallback := fallbackHandler{}
	h := gcshandler.New(gcshandler.Config{
		Bucket:   "acoshift",
		BasePath: "/",
		Fallback: &fallback,
	})

	r := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	assert.Equal(t, 200, w.Code)
	assert.False(t, fallback.called)
}

func TestHandlerNotFound(t *testing.T) {
	t.Run("WithFallback", func(t *testing.T) {
		fallback := fallbackHandler{}
		h := gcshandler.New(gcshandler.Config{
			Bucket:   "acoshift",
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
			Bucket:   "acoshift",
			BasePath: "/",
		})

		r := httptest.NewRequest(http.MethodGet, "http://localhost/not-exists-file", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		assert.Equal(t, 404, w.Code)
	})
}
