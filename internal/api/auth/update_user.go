package auth

import (
	"context"
	"log"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// UpdateUser updates user info by id
func (i *Implementation) UpdateUser(ctx context.Context, req *user_v1.UpdateUserRequest) error {
	err := i.authService.UpdateUser(ctx, converter.ToUserFromUserV1()) //*user_v1.User{}))
	if err != nil {
		return err
	}

	log.Printf("updated user with id: %d", req.GetId())

	return nil
}
