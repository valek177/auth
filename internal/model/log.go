package model

import "time"

type Record struct {
	ID        int64
	UserID    int64
	Action    string
	CreatedAt time.Time
}
