package model

import (
	"database/sql"
	"time"
)

type NewUser struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
}

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
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
