package model

import (
	"database/sql"
	"time"
)

// User contains user settings
type User struct {
	ID        int64         `db: "id"`
	UserInfo  UserInfo      `db: ""`
	CreatedAt time.Time     `db: "created_at"`
	UpdatedAt *sql.NullTime `db: "updated_at"`
}

// UserInfo contains user info
type UserInfo struct {
	Name  string `db: "name"`
	Email string `db: "email"`
	Role  string `db: "role"`
}
