package auth

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
)

// DeleteUser removes user
func (i *Implementation) DeleteUser(ctx context.Context, req *user_v1.DeleteUserRequest) (
	*emptypb.Empty, error,
) {
	err := i.authService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("deleted user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}
