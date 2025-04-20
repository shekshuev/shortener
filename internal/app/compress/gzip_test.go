package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipWriter_WriteAndClose(t *testing.T) {
	rr := httptest.NewRecorder()
	gz := NewGzipWriter(rr)
	gz.WriteHeader(http.StatusOK)

	data := []byte("hello, compressed world!")
	_, err := gz.Write(data)
	assert.NoError(t, err)

	err = gz.Close()
	assert.NoError(t, err)

	// проверим, что заголовок установлен
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// пытаемся распаковать, чтобы проверить содержимое
	gr, err := gzip.NewReader(bytes.NewReader(rr.Body.Bytes()))
	assert.NoError(t, err)

	uncompressed, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, data, uncompressed)
}

func TestGzipWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	gz := NewGzipWriter(rr)

	gz.WriteHeader(http.StatusTeapot)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	assert.Equal(t, http.StatusTeapot, rr.Code)
	_ = gz.Close()
}

func TestGzipReader_ReadAndClose(t *testing.T) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	original := []byte("this is a test")
	_, err := zw.Write(original)
	assert.NoError(t, err)
	err = zw.Close()
	assert.NoError(t, err)

	gr, err := NewGzipReader(io.NopCloser(bytes.NewReader(buf.Bytes())))
	assert.NoError(t, err)

	out, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, original, out)

	err = gr.Close()
	assert.NoError(t, err)
}

func TestGzipReader_InvalidInput(t *testing.T) {
	_, err := NewGzipReader(io.NopCloser(bytes.NewReader([]byte("not gzip"))))
	assert.Error(t, err)
}
