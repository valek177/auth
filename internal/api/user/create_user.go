package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// CreateUser creates new user with specified parameters
func (i *Implementation) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (
	*user_v1.CreateUserResponse, error,
) {
	err := validateCreateUser(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := i.userService.CreateUser(ctx, converter.ToNewUserFromNewUserV1(req))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &user_v1.CreateUserResponse{
		Id: id,
	}, nil
}
