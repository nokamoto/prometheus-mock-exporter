// Description:
// A simple Prometheus exporter that exposes a /metrics endpoint that prometheus can scrape.
//
// Usage:
//
//	go run cmd/prometheus-mock-exporter
//
// Environment variables:
//
//	PROMETHEUS_EXPORTER_ADDRESS - The address to listen on. Default is ":8080".
//	PROMETHEUS_EXPORTER_DEBUG - Enable debug mode. Set to any value to enable.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	envAddress        = "PROMETHEUS_EXPORTER_ADDRESS"
	envDefaultAddress = ":8080"

	envDebug = "PROMETHEUS_EXPORTER_DEBUG"
)

func run(ctx context.Context) {
	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "todo",
			Name:      "todo",
			Help:      "todo",
		},
		[]string{"todo"},
	)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down")
			return
		case <-time.After(time.Second):
			slog.Info("Incrementing counter")
			counter.WithLabelValues("foo").Inc()
		}
	}
}

func main() {
	address := os.Getenv(envAddress)
	if address == "" {
		address = envDefaultAddress
	}

	debug := os.Getenv(envDebug) != ""

	if debug {
		level := slog.LevelDebug
		slog.SetLogLoggerLevel(level)
		slog.Debug("Debug mode enabled", "level", level)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go run(ctx)

	slog.Info("Starting server", "address", address)
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: address}

	go func() {
		<-ctx.Done()
		slog.Info("Shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to gracefully shutdown server", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
