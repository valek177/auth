package auth

import (
	"context"
	"log"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/converter"
)

// GetUser returns info about user
func (i *Implementation) GetUser(ctx context.Context, req *user_v1.GetUserRequest) (
	*user_v1.GetUserResponse, error,
) {
	userObj, err := i.authService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n",
		userObj.ID, userObj.Name, userObj.Email, userObj.Role,
		userObj.CreatedAt, userObj.UpdatedAt)

	return &user_v1.GetUserResponse{
		User: converter.ToUserV1FromService(userObj),
	}, nil
}
