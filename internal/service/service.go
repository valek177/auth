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

// AuthService is interface for auth logic on service
type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService is interface for access logic on service
type AccessService interface {
	Check(ctx context.Context, accessToken string, endpoint string) (bool, error)
}

// ConsumerService is interface for consumer logic
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
