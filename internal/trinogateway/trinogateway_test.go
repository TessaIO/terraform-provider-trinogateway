// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trinogateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/client"
)

func TestBackendService(t *testing.T) {
	// Mock backends data
	mockBackends := []Backend{
		{
			Name:         "backend1",
			ProxyTo:      "http://trino1:8080",
			Active:       true,
			RoutingGroup: "default",
		},
		{
			Name:         "backend2",
			ProxyTo:      "http://trino2:8080",
			Active:       false,
			RoutingGroup: "batch",
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/entity/GATEWAY_BACKEND":
			if err := json.NewEncoder(w).Encode(mockBackends); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		case r.Method == "GET" && r.URL.Path == "/entity/GATEWAY_BACKEND/backend1":
			if err := json.NewEncoder(w).Encode(mockBackends[0]); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		case r.Method == "GET" && r.URL.Path == "/api/public/backends/backend1":
			if err := json.NewEncoder(w).Encode(mockBackends[0]); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		case r.Method == "POST" && r.URL.Path == "/entity":
			var req CreateBackendRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				fmt.Println("error while decoding response")
				return
			}
			backend := Backend{
				Name:         req.Name,
				ProxyTo:      req.ProxyTo,
				Active:       req.Active,
				RoutingGroup: req.RoutingGroup,
			}
			if err := json.NewEncoder(w).Encode(backend); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		case r.Method == "PUT" && r.URL.Path == "/entity/GATEWAY_BACKEND/backend1":
			var req UpdateBackendRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				fmt.Println("error while decoding response")
				return
			}
			backend := mockBackends[0]
			if req.Active != nil {
				backend.Active = *req.Active
			}
			if err := json.NewEncoder(w).Encode(backend); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		case r.Method == "POST" && r.URL.Path == "/gateway/backend/modify/delete":
			w.WriteHeader(http.StatusNoContent)

		case r.Method == "GET" && r.URL.Path == "/gateway/backend/state":
			states := []BackendState{
				{Name: "backend1", Active: true, Healthy: true},
				{Name: "backend2", Active: false, Healthy: false},
			}
			if err := json.NewEncoder(w).Encode(states); err != nil {
				fmt.Println("error while encoding response")
				return
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := client.NewClient(server.URL)
	svc := NewTrinoGatewayService(client)
	ctx := context.Background()

	t.Run("ListBackends", func(t *testing.T) {
		backends, err := svc.ListBackends(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(backends) != 2 {
			t.Errorf("expected 2 backends, got %d", len(backends))
		}
		if backends[0].Name != "backend1" {
			t.Errorf("expected backend1, got %s", backends[0].Name)
		}
	})

	t.Run("GetBackend", func(t *testing.T) {
		backend, err := svc.GetBackend(ctx, "backend1")
		if err != nil {
			t.Fatal(err)
		}
		if backend.Name != "backend1" {
			t.Errorf("expected backend1, got %s", backend.Name)
		}
		if !backend.Active {
			t.Error("expected backend to be active")
		}
	})

	t.Run("CreateBackend", func(t *testing.T) {
		req := CreateBackendRequest{
			Name:         "backend3",
			ProxyTo:      "http://trino3:8080",
			Active:       true,
			RoutingGroup: "adhoc",
		}
		backend, err := svc.CreateBackend(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if backend.Name != "backend3" {
			t.Errorf("expected backend3, got %s", backend.Name)
		}
	})

	t.Run("UpdateBackend", func(t *testing.T) {
		active := false
		req := UpdateBackendRequest{Active: &active}
		err := svc.UpdateBackend(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeleteBackend", func(t *testing.T) {
		err := svc.DeleteBackend(ctx, "backend1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ActivateBackend", func(t *testing.T) {
		err := svc.ActivateBackend(ctx, "backend1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DeactivateBackend", func(t *testing.T) {
		err := svc.DeactivateBackend(ctx, "backend1")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetBackendStates", func(t *testing.T) {
		states, err := svc.GetBackendStates(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(states) != 2 {
			t.Errorf("expected 2 states, got %d", len(states))
		}
		if !states[0].Healthy {
			t.Error("expected backend1 to be healthy")
		}
	})
}
