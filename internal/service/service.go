package service

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// UserService is interface for user logic on service
type UserService interface {
	CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

type AuthService interface{}

type AccessService interface{}

// ConsumerService is interface for consumer logic
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
