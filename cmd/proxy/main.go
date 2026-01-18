package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
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
	loader := proxy.NewConfigLoader()
	config, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.Debug {
		fmt.Fprintf(w, "Debug mode enabled\n")
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
		fmt.Fprintf(w, "Starting proxy server on %s\n", config.ListenAddr)
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
