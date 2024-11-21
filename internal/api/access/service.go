package access

import (
	"github.com/valek177/auth/grpc/pkg/access_v1"
	"github.com/valek177/auth/internal/service"
)

// Implementation struct contains server
type Implementation struct {
	access_v1.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewImplementation returns implementation object
func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
