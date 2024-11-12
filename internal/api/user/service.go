package user

import (
	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/service"
)

// Implementation struct contains server
type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation returns implementation object
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
