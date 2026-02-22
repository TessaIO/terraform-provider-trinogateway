// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trinogateway

import (
	"context"
	"fmt"

	"github.com/TessaIO/terraform-provider-trinogateway/internal/client"
)

// BackendService provides methods for managing Trino Gateway backends.
type TrinoGateway struct {
	client *client.Client
}

// NewBackendService creates a new BackendService.
func NewTrinoGatewayService(client *client.Client) *TrinoGateway {
	return &TrinoGateway{client: client}
}

// ListBackends retrieves all backends.
func (t *TrinoGateway) ListBackends(ctx context.Context) ([]Backend, error) {
	resp, err := t.client.Get(ctx, "/entity/GATEWAY_BACKEND")
	if err != nil {
		return nil, fmt.Errorf("failed to list backends: %w", err)
	}

	var backends []Backend
	if err := client.DecodeResponse(resp, &backends); err != nil {
		return nil, err
	}

	return backends, nil
}

// GetBackend retrieves a specific backend by name.
func (t *TrinoGateway) GetBackend(ctx context.Context, name string) (*Backend, error) {
	path := fmt.Sprintf("/api/public/backends/%s", name)
	resp, err := t.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get backend %s: %w", name, err)
	}

	var backend Backend
	if err := client.DecodeResponse(resp, &backend); err != nil {
		return nil, err
	}

	return &backend, nil
}

// CreateBackend creates a new backend.
func (t *TrinoGateway) CreateBackend(ctx context.Context, req CreateBackendRequest) (*Backend, error) {
	_, err := t.client.Post(ctx, "/entity?entityType=GATEWAY_BACKEND", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend: %w", err)
	}

	// API Response from the Gateway doesn't return anything so we construct the object here
	backend := Backend(req)
	return &backend, nil
}

// UpdateBackend updates an existing backend.
func (t *TrinoGateway) UpdateBackend(ctx context.Context, req UpdateBackendRequest) error {
	path := "/entity?entityType=GATEWAY_BACKEND"
	_, err := t.client.Post(ctx, path, req)
	if err != nil {
		return fmt.Errorf("failed to update backend %s: %w", req.Name, err)
	}

	return nil
}

// DeleteBackend deletes a backend.
func (t *TrinoGateway) DeleteBackend(ctx context.Context, name string) error {
	path := "/gateway/backend/modify/delete"
	_, err := t.client.Post(ctx, path, name)
	if err != nil {
		return fmt.Errorf("failed to delete backend %s: %w", name, err)
	}

	return nil
}

// ActivateBackend activates a backend.
func (t *TrinoGateway) ActivateBackend(ctx context.Context, name string) error {
	active := true
	req := UpdateBackendRequest{Active: &active, Name: name}
	err := t.UpdateBackend(ctx, req)
	return err
}

// DeactivateBackend deactivates a backend.
func (t *TrinoGateway) DeactivateBackend(ctx context.Context, name string) error {
	active := false
	req := UpdateBackendRequest{Active: &active, Name: name}
	err := t.UpdateBackend(ctx, req)
	return err
}

// GetBackendStates retrieves the state of all backends.
func (t *TrinoGateway) GetBackendStates(ctx context.Context) ([]BackendState, error) {
	resp, err := t.client.Get(ctx, "/gateway/backend/state")
	if err != nil {
		return nil, fmt.Errorf("failed to get backend states: %w", err)
	}

	var states []BackendState
	if err := client.DecodeResponse(resp, &states); err != nil {
		return nil, err
	}

	return states, nil
}
