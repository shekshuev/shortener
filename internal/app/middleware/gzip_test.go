package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func gzipCompress(data string) []byte {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, _ = w.Write([]byte(data))
	_ = w.Close()
	return buf.Bytes()
}

func TestGzipCompressor_RequestWithGzipBody(t *testing.T) {
	data := "Hello, world!"
	body := gzipCompress(data)

	var received string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		received = string(b)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	GzipCompressor(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, data, received)
}

func TestGzipCompressor_ResponseWithGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("compressed response"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	GzipCompressor(handler).ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	gr, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err)
	out, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, "compressed response", string(out))
}

func TestGzipCompressor_NoGzipHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("no compression"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	GzipCompressor(handler).ServeHTTP(w, req)

	resp := w.Result()
	assert.Empty(t, resp.Header.Get("Content-Encoding"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "no compression", string(body))
}

func TestGzipCompressor_InvalidGzipBody(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.FailNow()
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not gzipped data"))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	GzipCompressor(handler).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
