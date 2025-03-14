package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// GzipWriter представляет структуру для сжатия HTTP-ответа.
type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewGzipWriter создает новый экземпляр GzipWriter
func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header возвращает заголовки оригинального writer'а.
func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

// Write сжимает и записывает данные в ответ.
func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader устанавливает код статуса ответа и указывает, что содержимое сжато gzip.
func (c *GzipWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close завершает работу и закрывает gzip writer.
func (c *GzipWriter) Close() error {
	return c.zw.Close()
}

// GzipReader оборачивает io.ReadCloser и выполняет декомпрессию данных, сжатых gzip.
type GzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewGzipReader инициализирует новый экземпляр GzipReader.
func NewGzipReader(r io.ReadCloser) (*GzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read читает декомпрессированные данные в p.
func (c GzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает оригинальный reader и gzip reader.
func (c *GzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
