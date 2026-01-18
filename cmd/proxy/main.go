package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/braw-dev/memex/internal/proxy"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	w io.Writer,
	args []string,
	getenv func(string) string,
) error {
	// Initialize config
	loader := proxy.NewConfigLoader(getenv)
	config, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := setupLogger(config.Log); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	// Create handler using NewServer (returns http.Handler)
	handler := proxy.NewServer(config)

	// Initialize http.Server
	srv := &http.Server{
		Addr:         config.ListenAddr,
		Handler:      handler,
		ReadTimeout:  config.UpstreamTimeout, // Approximate
		WriteTimeout: config.UpstreamTimeout, // Approximate
		IdleTimeout:  config.IdleTimeout,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		slog.Info("Starting proxy server", "addr", config.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for signal or error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed: %w", err)
	case sig := <-sigChan:
		fmt.Fprintf(w, "\nReceived signal %v, shutting down...\n", sig)
	case <-ctx.Done():
		fmt.Fprintf(w, "\nContext cancelled, shutting down...\n")
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	fmt.Fprintf(w, "Server stopped\n")
	return nil
}

func setupLogger(cfg proxy.LogConfig) error {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var w io.Writer = os.Stderr
	if cfg.Path != "" && cfg.Path != "-" && cfg.Path != "stderr" {
		if cfg.Path == "stdout" {
			w = os.Stdout
		} else {
			f, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			w = f
		}
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(w, opts)
	} else {
		handler = slog.NewTextHandler(w, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}
