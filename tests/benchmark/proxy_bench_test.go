package benchmark

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/braw-dev/memex/internal/proxy"
)

func BenchmarkProxyPassthrough(b *testing.B) {
	// Setup upstream
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	// Setup proxy
	config := &proxy.ProxyConfig{
		ListenAddr:      ":0",
		UpstreamTimeout: 5 * time.Second,
		IdleTimeout:     5 * time.Second,
		FlushInterval:   0,
		Log:             proxy.LogConfig{Level: "error"}, // Minimal logging
	}
	handler := proxy.NewServer(config)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	proxyURL, _ := url.Parse(proxyServer.URL)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := client.Get(upstream.URL)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}
