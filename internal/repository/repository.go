package repository

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// AuthRepository is interface for user logic
type AuthRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error
}
