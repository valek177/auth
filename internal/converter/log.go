package converter

import (
	"time"

	"github.com/valek177/auth/internal/model"
)

// ToRecordRepoFromService converts params to Record model
func ToRecordRepoFromService(userID int64, action string) *model.Record {
	return &model.Record{
		UserID:    userID,
		CreatedAt: time.Now(),
		Action:    action,
	}
}
