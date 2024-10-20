package auth

import (
	"context"
	"log"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// CreateUser creates new user with specified parameters
func (i *Implementation) CreateUser(ctx context.Context, req *user_v1.CreateUserRequest) (
	*user_v1.CreateUserResponse, error,
) {
	id, err := i.authService.CreateUser(ctx, converter.ToUserFromUserV1(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &user_v1.CreateUserResponse{
		Id: id,
	}, nil
}
