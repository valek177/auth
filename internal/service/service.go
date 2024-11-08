package service

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// AuthService is interface for user logic on service
type AuthService interface {
	CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
