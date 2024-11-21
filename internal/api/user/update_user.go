package user

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// UpdateUser updates user info by id
func (i *Implementation) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (
	*emptypb.Empty, error,
) {
	err := validateUpdateUser(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = i.userService.UpdateUser(ctx, converter.ToUpdateUserInfoFromV1(req))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &emptypb.Empty{}, nil
}
