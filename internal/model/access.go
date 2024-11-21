package model

// EndpointAccessRule is a model for access rule of endpoint
type EndpointAccessRule struct {
	Endpoint string
	Roles    []string
}
