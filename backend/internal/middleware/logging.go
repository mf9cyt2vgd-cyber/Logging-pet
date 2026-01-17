package middleware

import (
	"litestar/pkg/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg := logger.NewLogger()
		start := time.Now()
		logg.Info("request",
			"Host:", r.Host,
			"Method:", r.Method,
			"Path:", r.URL.Path)
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logg.Info("request complete", "duration", duration)
	})
}
