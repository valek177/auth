package repository

import (
	"context"

	"github.com/valek177/auth/internal/model"
)

// AuthRepository is interface for user logic
type AuthRepository interface {
	CreateUser(ctx context.Context, newUser *model.NewUser) (int64, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	UpdateUser(ctx context.Context, updateUserInfo *model.UpdateUserInfo) error
	DeleteUser(ctx context.Context, id int64) error
}

type LogRepository interface {
	CreateRecord(ctx context.Context, record *model.Record) (int64, error)
}
