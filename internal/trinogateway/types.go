package trinogateway

import "time"

// Backend represents a Trino backend configuration
type Backend struct {
	Name         string `json:"name"`
	ProxyTo      string `json:"proxyTo"`
	Active       bool   `json:"active"`
	RoutingGroup string `json:"routingGroup,omitempty"`
	ExternalURL  string `json:"externalUrl,omitempty"`
	LastHealthy  string `json:"lastHealthy,omitempty"`
}

// BackendState represents the state of a backend
type BackendState struct {
	Name         string    `json:"name"`
	Active       bool      `json:"active"`
	Healthy      bool      `json:"healthy"`
	LastHealthy  time.Time `json:"lastHealthy,omitempty"`
	RoutingGroup string    `json:"routingGroup,omitempty"`
}

// QueryHistory represents a query in the history
type QueryHistory struct {
	QueryID    string    `json:"queryId"`
	User       string    `json:"user"`
	Source     string    `json:"source"`
	BackendURL string    `json:"backendUrl"`
	Created    time.Time `json:"created"`
	QueryText  string    `json:"queryText,omitempty"`
}

// RoutingGroup represents a routing group configuration
type RoutingGroup struct {
	Name     string `json:"name"`
	Selector string `json:"selector,omitempty"`
}

// GatewayBackendConfiguration represents the gateway configuration
type GatewayBackendConfiguration struct {
	Name         string `json:"name"`
	ProxyTo      string `json:"proxyTo"`
	Active       bool   `json:"active"`
	RoutingGroup string `json:"routingGroup"`
	ExternalURL  string `json:"externalUrl"`
}

// CreateBackendRequest represents the request to create a new backend
type CreateBackendRequest struct {
	Name         string `json:"name"`
	ProxyTo      string `json:"proxyTo"`
	Active       bool   `json:"active"`
	RoutingGroup string `json:"routingGroup,omitempty"`
	ExternalURL  string `json:"externalUrl,omitempty"`
}

// UpdateBackendRequest represents the request to update a backend
type UpdateBackendRequest struct {
	Name         string  `json:"name,omitempty"`
	Active       *bool   `json:"active,omitempty"`
	RoutingGroup *string `json:"routingGroup,omitempty"`
	ExternalURL  *string `json:"externalUrl,omitempty"`
	ProxyTo      *string `json:"proxyTo,omitempty"`
}
