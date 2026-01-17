package main

import (
	"context"
	"errors"
	"litestar/internal/handlers"
	"litestar/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(middleware.RequestsTotal, middleware.RequestsDuration)
}
func main() {
	done := make(chan struct{})
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/status/", handlers.StatusHandler)
	handler := middleware.LoggingMiddleware(
		middleware.MetricsMiddleware(
			mux,
		),
	)
	srv := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	slog.Info("Starting server...")
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("fatal error:", err)
			return
		}
	}()
	go func() {
		defer close(done)
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		slog.Warn("closing server...shutdown...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	<-done
}
