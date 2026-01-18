package proxy

import (
	"os"
	"testing"
)

func TestDefaultConfigLoader_Load(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("MEMEX_PROXY_LOG_LEVEL", "debug")
	os.Setenv("MEMEX_PROXY_LOG_FORMAT", "json")
	defer os.Unsetenv("MEMEX_PROXY_LOG_LEVEL")
	defer os.Unsetenv("MEMEX_PROXY_LOG_FORMAT")

	loader := NewConfigLoader()
	config, err := loader.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if config.Log.Level != "debug" {
		t.Errorf("expected level 'debug', got %s", config.Log.Level)
	}
	if config.Log.Format != "json" {
		t.Errorf("expected format 'json', got %s", config.Log.Format)
	}
}

func TestDefaults(t *testing.T) {
	// Ensure clean env (only relevant keys)
	os.Unsetenv("MEMEX_PROXY_LOG_LEVEL")
	os.Unsetenv("MEMEX_PROXY_LOG_FORMAT")
	os.Unsetenv("MEMEX_PROXY_LOG_PATH")

	loader := NewConfigLoader()
	config, err := loader.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if config.Log.Level != "info" {
		t.Errorf("expected default level 'info', got %s", config.Log.Level)
	}
	if config.Log.Format != "text" {
		t.Errorf("expected default format 'text', got %s", config.Log.Format)
	}
}
