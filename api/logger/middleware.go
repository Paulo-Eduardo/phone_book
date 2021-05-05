package logger

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
  http.ResponseWriter
  statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
  return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
  lrw.statusCode = code
  lrw.ResponseWriter.WriteHeader(code)
}


func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    t := time.Now()

		lrw := NewLoggingResponseWriter(w)
		defer func() {
			log.Printf("(HTTP %v) %v %v in %v", lrw.statusCode, r.Method, r.RequestURI, time.Since(t))
		}()

		handler.ServeHTTP(lrw, r)
	})
}
