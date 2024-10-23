package auth

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// UpdateUser updates user info by id
func (i *Implementation) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) (
	*emptypb.Empty, error,
) {
	isUpdated := false
	if req.GetName() != nil {
		isUpdated = true
	}
	if req.GetRole().String() != "" {
		isUpdated = true
	}
	if !isUpdated {
		return &emptypb.Empty{}, nil
	}
	err := i.authService.UpdateUser(ctx, converter.ToUpdateUserInfoFromV1(req))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
