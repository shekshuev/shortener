package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLogger_WritesResponseCorrectly(t *testing.T) {
	handlerCalled := false

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("Hello, logs!"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	RequestLogger(h).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.True(t, handlerCalled, "Handler should be called")
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	assert.Equal(t, "Hello, logs!", string(body))
}

func TestRequestLogger_EmptyResponse(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/empty", nil)
	rec := httptest.NewRecorder()

	RequestLogger(h).ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, 0, len(body))
}
