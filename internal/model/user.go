package model

import (
	"database/sql"
	"time"
)

// NewUser is a model for created user
type NewUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

// UpdateUserInfo is a model for updated params of user
type UpdateUserInfo struct {
	ID   int64
	Name *string
	Role *string
}

// User contains user settings
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
