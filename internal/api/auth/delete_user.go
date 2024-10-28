package auth

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
)

// DeleteUser removes user
func (i *Implementation) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (
	*emptypb.Empty, error,
) {
	err := validateDeleteUser(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = i.authService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
