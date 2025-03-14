package middleware

import (
	"net/http"
	"time"

	"github.com/shekshuev/shortener/internal/app/logger"
	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write записывает данные в ответ и обновляет размер ответа.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader записывает HTTP-статус в ответ и фиксирует его в responseData.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// RequestLogger - middleware для логирования HTTP-запросов и ответов.
func RequestLogger(h http.Handler) http.Handler {
	log := logger.NewLogger()
	logFn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		log.Log.Info("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("uri", r.URL.Path),
			zap.String("duration", duration.String()),
		)
		log.Log.Info("got outgoing HTTP response",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
	return http.HandlerFunc(logFn)
}
