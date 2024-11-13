package access

import (
	"github.com/pkg/errors"

	"github.com/valek177/auth/grpc/pkg/access_v1"
)

func validateCheck(req *access_v1.CheckRequest) error {
	if req == nil {
		return errors.New("unable to check access: empty request")
	}

	return nil
}
