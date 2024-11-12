package auth

import (
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	// refresh token? access?
}

// NewService creates new service with settings
func NewService(
	userRepository repository.UserRepository,
) service.AuthService {
	return &serv{
		userRepository: userRepository,
	}
}
