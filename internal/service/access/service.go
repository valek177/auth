package access

import (
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
	"github.com/valek177/auth/internal/utils"
)

type serv struct {
	accessRepository repository.AccessRepository
	tokenAccess      utils.Token
}

// NewService creates new service with settings
func NewService(
	accessRepository repository.AccessRepository,
	tokenAccess utils.Token,
) service.AccessService {
	return &serv{
		accessRepository: accessRepository,
		tokenAccess:      tokenAccess,
	}
}
