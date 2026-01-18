package proxy

import (
	"testing"
)

func TestDefaultConfigLoader_Load(t *testing.T) {
	// Mock getenv
	env := map[string]string{
		"MEMEX_PROXY_LOG_LEVEL":  "debug",
		"MEMEX_PROXY_LOG_FORMAT": "json",
	}
	getenv := func(key string) string {
		return env[key]
	}

	loader := NewConfigLoader(getenv)
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
	// Empty getenv
	getenv := func(key string) string { return "" }

	loader := NewConfigLoader(getenv)
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
