package integration

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/braw-dev/memex/internal/proxy"
)

func TestHTTPPassthrough(t *testing.T) {
	// 1. Setup Mock Upstream Server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Test") != "true" {
			t.Errorf("Expected X-Test header")
		}
		// Verify Proxy- headers are removed
		if r.Header.Get("Proxy-Authorization") != "" {
			t.Errorf("Expected Proxy-Authorization to be removed")
		}
		
		w.Header().Set("X-Response", "ok")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer upstream.Close()

	// 2. Setup Proxy
	config := &proxy.ProxyConfig{
		ListenAddr:      ":0", // Random port
		UpstreamTimeout: 5 * time.Second,
		IdleTimeout:     5 * time.Second,
		FlushInterval:   0,
		Debug:           true,
	}
	
	handler := proxy.NewServer(config)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	// 3. Make Request through Proxy
	proxyURL, _ := url.Parse(proxyServer.URL)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	req, _ := http.NewRequest("GET", upstream.URL+"/test", nil)
	req.Header.Set("X-Test", "true")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// 4. Verify Response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Response") != "ok" {
		t.Errorf("Expected X-Response header")
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "hello world" {
		t.Errorf("Expected body 'hello world', got '%s'", string(body))
	}
}

func TestHealthz(t *testing.T) {
	config := &proxy.ProxyConfig{
		ListenAddr: ":0",
		Debug:      false,
	}
	handler := proxy.NewServer(config)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	resp, err := http.Get(proxyServer.URL + "/healthz")
	if err != nil {
		t.Fatalf("Health check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", string(body))
	}
}

func TestProxyHeaders(t *testing.T) {
    // Test that Proxy-Authorization header is stripped
    upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Proxy-Authorization") != "" {
            t.Errorf("Proxy-Authorization header was not stripped")
        }
        w.WriteHeader(http.StatusOK)
    }))
    defer upstream.Close()

    config := &proxy.ProxyConfig{
        ListenAddr: ":0",
        Debug:      true,
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

    req, _ := http.NewRequest("GET", upstream.URL, nil)
    req.Header.Set("Proxy-Authorization", "Basic user:pass") // Manually set header to test stripping
    
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }
    resp.Body.Close()
}
