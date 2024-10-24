package converter

import (
	"time"

	"github.com/valek177/auth/internal/model"
)

func ToRecordRepoFromService(userId int64, action string) *model.Record {
	return &model.Record{
		UserID:    userId,
		CreatedAt: time.Now(),
		Action:    action,
	}
}