package auth

import (
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
	"github.com/valek177/auth/internal/utils"
)

type serv struct {
	userRepository repository.UserRepository
	tokenRefresh   utils.Token
	tokenAccess    utils.Token
}

// NewService creates new service with settings
func NewService(
	userRepository repository.UserRepository,
	tokenRefresh utils.Token,
	tokenAccess utils.Token,
) service.AuthService {
	return &serv{
		userRepository: userRepository,
		tokenRefresh:   tokenRefresh,
		tokenAccess:    tokenAccess,
	}
}
