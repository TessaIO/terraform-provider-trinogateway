package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		opts        []Option
		expectError bool
	}{
		{
			name:        "valid URL without auth",
			baseURL:     "http://localhost:8080",
			opts:        nil,
			expectError: false,
		},
		{
			name:        "valid URL with auth",
			baseURL:     "http://localhost:8080",
			opts:        []Option{WithAuth("user", "pass")},
			expectError: false,
		},
		{
			name:        "empty URL",
			baseURL:     "",
			opts:        nil,
			expectError: true,
		},
		{
			name:        "invalid URL",
			baseURL:     "://invalid",
			opts:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.opts...)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client == nil {
					t.Error("expected client but got nil")
				}
			}
		})
	}
}

func TestClientWithOptions(t *testing.T) {
	t.Run("WithAuth", func(t *testing.T) {
		client, err := NewClient("http://localhost:8080", WithAuth("testuser", "testpass"))
		if err != nil {
			t.Fatal(err)
		}
		if client.auth == nil {
			t.Fatal("auth should not be nil")
		}
		if client.auth.Username != "testuser" || client.auth.Password != "testpass" {
			t.Error("auth credentials not set correctly")
		}
	})

	t.Run("WithTimeout", func(t *testing.T) {
		timeout := 45 * time.Second
		client, err := NewClient("http://localhost:8080", WithTimeout(timeout))
		if err != nil {
			t.Fatal(err)
		}
		if client.httpClient.Timeout != timeout {
			t.Errorf("expected timeout %v, got %v", timeout, client.httpClient.Timeout)
		}
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		customClient := &http.Client{Timeout: 60 * time.Second}
		client, err := NewClient("http://localhost:8080", WithHTTPClient(customClient))
		if err != nil {
			t.Fatal(err)
		}
		if client.httpClient != customClient {
			t.Error("custom HTTP client not set")
		}
	})
}

func TestClientRequests(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		username, password, ok := r.BasicAuth()
		if ok && username == "testuser" && password == "testpass" {
			w.Header().Set("X-Auth", "verified")
		}

		// Return different responses based on path
		switch r.URL.Path {
		case "/test":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "success"})
		case "/error":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		default:
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"path": r.URL.Path})
		}
	}))
	defer server.Close()

	t.Run("GET request", func(t *testing.T) {
		client, _ := NewClient(server.URL)
		ctx := context.Background()

		resp, err := client.Get(ctx, "/test")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GET request with auth", func(t *testing.T) {
		client, _ := NewClient(server.URL, WithAuth("testuser", "testpass"))
		ctx := context.Background()

		resp, err := client.Get(ctx, "/test")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.Header.Get("X-Auth") != "verified" {
			t.Error("authentication not applied")
		}
	})

	t.Run("POST request", func(t *testing.T) {
		client, _ := NewClient(server.URL)
		ctx := context.Background()

		body := map[string]string{"key": "value"}
		resp, err := client.Post(ctx, "/test", body)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("HTTP error", func(t *testing.T) {
		client, _ := NewClient(server.URL)
		ctx := context.Background()

		_, err := client.Get(ctx, "/error")
		if err == nil {
			t.Fatal("expected error but got none")
		}

		httpErr, ok := err.(*HTTPError)
		if !ok {
			t.Fatal("expected HTTPError type")
		}
		if httpErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", httpErr.StatusCode)
		}
	})
}

func TestHTTPError(t *testing.T) {
	err := &HTTPError{
		StatusCode: 404,
		Status:     "Not Found",
		Body:       "resource not found",
	}

	expected := "HTTP 404: Not Found - resource not found"
	if err.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, err.Error())
	}
}

func TestDecodeResponse(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]string{"name": "test"})
		}))
		defer server.Close()

		client, _ := NewClient(server.URL)
		resp, _ := client.Get(context.Background(), "/")

		var result map[string]string
		err := DecodeResponse(resp, &result)
		if err != nil {
			t.Fatal(err)
		}
		if result["name"] != "test" {
			t.Errorf("expected name 'test', got %q", result["name"])
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client, _ := NewClient(server.URL)
		resp, _ := client.Get(context.Background(), "/")

		var result map[string]string
		err := DecodeResponse(resp, &result)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})
}
