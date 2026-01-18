package proxy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
)

// LogConfig represents the logging configuration
type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
	Path   string `koanf:"path"`
}

// ProxyConfig represents the proxy server configuration
type ProxyConfig struct {
	ListenAddr      string        `koanf:"listen"`
	UpstreamTimeout time.Duration `koanf:"upstream_timeout"`
	IdleTimeout     time.Duration `koanf:"idle_timeout"`
	FlushInterval   time.Duration `koanf:"flush_interval"`
	Log             LogConfig     `koanf:"log"`
}

// ConfigLoader loads configuration from various sources
type ConfigLoader interface {
	Load() (*ProxyConfig, error)
}

// DefaultConfigLoader implements ConfigLoader using koanf
type DefaultConfigLoader struct {
	getenv func(string) string
}

// NewConfigLoader creates a new DefaultConfigLoader
func NewConfigLoader(getenv func(string) string) *DefaultConfigLoader {
	return &DefaultConfigLoader{
		getenv: getenv,
	}
}

// Load loads configuration from files and environment
func (l *DefaultConfigLoader) Load() (*ProxyConfig, error) {
	k := koanf.New(".")

	// Default values
	defaults := map[string]interface{}{
		"proxy": map[string]interface{}{
			"listen":           ":8080",
			"upstream_timeout": "60s",
			"idle_timeout":     "90s",
			"flush_interval":   "0s",
			"log": map[string]interface{}{
				"level":  "info",
				"format": "text",
				"path":   "stderr",
			},
		},
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

	// Load environment variables
	if err := k.Load(NewEnvProvider("MEMEX", l.getenv), nil); err != nil {
		return nil, fmt.Errorf("failed to load env vars: %w", err)
	}

	// Unmarshal into struct
	config := &ProxyConfig{}
	if err := k.Unmarshal("proxy", config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// NewEnvProvider creates a koanf.Provider that loads environment variables
// using a custom getenv function.
func NewEnvProvider(prefix string, getenv func(string) string) koanf.Provider {
	return &envProvider{
		prefix: prefix,
		getenv: getenv,
	}
}

type envProvider struct {
	prefix string
	getenv func(string) string
}

func (p *envProvider) ReadBytes() ([]byte, error) {
	return nil, nil
}

func (p *envProvider) Read() (map[string]interface{}, error) {
	// Since we can't iterate environment variables using just getenv,
	// we have to rely on known configuration keys.
	// This provider explicitly looks up keys that match our configuration structure.
	// This is a trade-off to allow testability without os.Setenv.

	m := make(map[string]interface{})
	p.lookup(m, "proxy.listen", "PROXY_LISTEN")
	p.lookup(m, "proxy.upstream_timeout", "PROXY_UPSTREAM_TIMEOUT")
	p.lookup(m, "proxy.idle_timeout", "PROXY_IDLE_TIMEOUT")
	p.lookup(m, "proxy.flush_interval", "PROXY_FLUSH_INTERVAL")
	p.lookup(m, "proxy.log.level", "PROXY_LOG_LEVEL")
	p.lookup(m, "proxy.log.format", "PROXY_LOG_FORMAT")
	p.lookup(m, "proxy.log.path", "PROXY_LOG_PATH")

	return m, nil
}

func (p *envProvider) lookup(m map[string]interface{}, configKey, envSuffix string) {
	envKey := p.prefix + "_" + envSuffix
	if val := p.getenv(envKey); val != "" {
		// Koanf requires keys to be nested maps for nested structs when loading from a map
		// We need to split the key by "." and create nested maps
		keys := strings.Split(configKey, ".")
		current := m
		for i, k := range keys {
			if i == len(keys)-1 {
				current[k] = val
			} else {
				if _, ok := current[k]; !ok {
					current[k] = make(map[string]interface{})
				}
				if next, ok := current[k].(map[string]interface{}); ok {
					current = next
				} else {
					// Should not happen if keys are distinct
					break
				}
			}
		}
	}
}

func (p *envProvider) Watch(cb func(event interface{}, err error)) error {
	return nil
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
