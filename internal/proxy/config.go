package proxy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
)

// ProxyConfig represents the proxy server configuration
type ProxyConfig struct {
	ListenAddr      string        `koanf:"listen"`
	UpstreamTimeout time.Duration `koanf:"upstream_timeout"`
	IdleTimeout     time.Duration `koanf:"idle_timeout"`
	FlushInterval   time.Duration `koanf:"flush_interval"`
	Debug           bool          `koanf:"debug"`
}

// ConfigLoader loads configuration from various sources
type ConfigLoader interface {
	Load() (*ProxyConfig, error)
}

// DefaultConfigLoader implements ConfigLoader using koanf
type DefaultConfigLoader struct{}

// NewConfigLoader creates a new DefaultConfigLoader
func NewConfigLoader() *DefaultConfigLoader {
	return &DefaultConfigLoader{}
}

// Load loads configuration from files and environment
func (l *DefaultConfigLoader) Load() (*ProxyConfig, error) {
	k := koanf.New(".")

	// Default values
	defaults := map[string]interface{}{
		"proxy.listen":           ":8080",
		"proxy.upstream_timeout": "60s",
		"proxy.idle_timeout":     "90s",
		"proxy.flush_interval":   "0s",
		"proxy.debug":            false,
	}
	if err := k.Load(mapProvider(defaults), nil); err != nil {
		return nil, fmt.Errorf("failed to load defaults: %w", err)
	}

	// Config file priority
	configFiles := []string{
		"memex.yml", ".memex.yml", ".config/memex.yml",
		"memex.yaml", ".memex.yaml", ".config/memex.yaml",
		"memex.toml", ".memex.toml", ".config/memex.toml",
	}

	// Try to find the first config file that exists
	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			// Found a file, load it
			var parser koanf.Parser
			ext := filepath.Ext(file)
			if ext == ".toml" {
				parser = toml.Parser()
			} else {
				parser = yaml.Parser()
			}

			// Use os.DirFS(".") to load from current directory
			if err := k.Load(fs.Provider(os.DirFS("."), file), parser); err != nil {
				return nil, fmt.Errorf("failed to load config file %s: %w", file, err)
			}
			break
		}
	}

	// Unmarshal into struct
	config := &ProxyConfig{}
	if err := k.Unmarshal("proxy", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// mapProvider is a simple provider for a map
func mapProvider(m map[string]interface{}) koanf.Provider {
	return &mp{m: m}
}

type mp struct {
	m map[string]interface{}
}

func (p *mp) Read() (map[string]interface{}, error) {
	return p.m, nil
}

func (p *mp) ReadBytes() ([]byte, error) {
	return nil, nil
}

func (p *mp) Watch(cb func(event interface{}, err error)) error {
	return nil
}
