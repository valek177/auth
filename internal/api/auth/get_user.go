package auth

import (
	"context"

	"github.com/pkg/errors"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// GetUser returns info about user
func (i *Implementation) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (
	*user_v1.GetUserResponse, error,
) {
	err := validateGetUser(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userObj, err := i.authService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &user_v1.GetUserResponse{
		User: converter.ToUserV1FromService(userObj),
	}, nil
}
