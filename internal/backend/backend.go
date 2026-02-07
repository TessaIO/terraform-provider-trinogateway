package backend

import (
	"context"
	"fmt"

	"github.com/TessaIO/terraform-provider-trino-gateway/internal/client"
)

// BackendService provides methods for managing Trino Gateway backends
type BackendService struct {
	client *client.Client
}

// NewBackendService creates a new BackendService
func NewBackendService(client *client.Client) *BackendService {
	return &BackendService{client: client}
}

// ListBackends retrieves all backends
func (s *BackendService) ListBackends(ctx context.Context) ([]Backend, error) {
	resp, err := s.client.Get(ctx, "/entity/GATEWAY_BACKEND")
	if err != nil {
		return nil, fmt.Errorf("failed to list backends: %w", err)
	}

	var backends []Backend
	if err := client.DecodeResponse(resp, &backends); err != nil {
		return nil, err
	}

	return backends, nil
}

// GetBackend retrieves a specific backend by name
func (s *BackendService) GetBackend(ctx context.Context, name string) (*Backend, error) {
	path := fmt.Sprintf("/entity/GATEWAY_BACKEND/%s", name)
	resp, err := s.client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get backend %s: %w", name, err)
	}

	var backend Backend
	if err := client.DecodeResponse(resp, &backend); err != nil {
		return nil, err
	}

	return &backend, nil
}

// CreateBackend creates a new backend
func (s *BackendService) CreateBackend(ctx context.Context, req CreateBackendRequest) (*Backend, error) {
	resp, err := s.client.Post(ctx, "/entity/GATEWAY_BACKEND", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend: %w", err)
	}

	var backend Backend
	if err := client.DecodeResponse(resp, &backend); err != nil {
		return nil, err
	}

	return &backend, nil
}

// UpdateBackend updates an existing backend
func (s *BackendService) UpdateBackend(ctx context.Context, name string, req UpdateBackendRequest) (*Backend, error) {
	path := fmt.Sprintf("/entity/GATEWAY_BACKEND/%s", name)
	resp, err := s.client.Put(ctx, path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update backend %s: %w", name, err)
	}

	var backend Backend
	if err := client.DecodeResponse(resp, &backend); err != nil {
		return nil, err
	}

	return &backend, nil
}

// DeleteBackend deletes a backend
func (s *BackendService) DeleteBackend(ctx context.Context, name string) error {
	path := fmt.Sprintf("/entity/GATEWAY_BACKEND/%s", name)
	_, err := s.client.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete backend %s: %w", name, err)
	}

	return nil
}

// ActivateBackend activates a backend
func (s *BackendService) ActivateBackend(ctx context.Context, name string) error {
	active := true
	req := UpdateBackendRequest{Active: &active}
	_, err := s.UpdateBackend(ctx, name, req)
	return err
}

// DeactivateBackend deactivates a backend
func (s *BackendService) DeactivateBackend(ctx context.Context, name string) error {
	active := false
	req := UpdateBackendRequest{Active: &active}
	_, err := s.UpdateBackend(ctx, name, req)
	return err
}

// GetBackendStates retrieves the state of all backends
func (s *BackendService) GetBackendStates(ctx context.Context) ([]BackendState, error) {
	resp, err := s.client.Get(ctx, "/gateway/backend/state")
	if err != nil {
		return nil, fmt.Errorf("failed to get backend states: %w", err)
	}

	var states []BackendState
	if err := client.DecodeResponse(resp, &states); err != nil {
		return nil, err
	}

	return states, nil
}
