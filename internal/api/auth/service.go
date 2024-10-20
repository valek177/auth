package auth

import (
	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/service"
)

// Implementation struct contains server
type Implementation struct {
	user_v1.UnimplementedUserV1Server
	authService service.AuthService
}

// NewImplementation returns implementation object
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
