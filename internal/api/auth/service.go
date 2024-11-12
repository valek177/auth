package auth

import (
	"github.com/valek177/auth/grpc/pkg/auth_v1"
	"github.com/valek177/auth/internal/service"
)

// Implementation struct contains server
type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation returns implementation object
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
