package access

import (
	"github.com/valek177/auth/internal/service"
)

type serv struct {
	// token provider utils ?
}

// NewService creates new service with settings
func NewService() service.AccessService {
	return &serv{}
}
