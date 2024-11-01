package converter

import (
	"database/sql"
	"time"

	"github.com/valek177/auth/internal/model"
	modelRepo "github.com/valek177/auth/internal/repository/redis/model"
)

// ToRedisRepoFromUser converts user from service to repo
func ToRedisRepoFromUser(user *model.User) *modelRepo.UserRedis {
	if user == nil {
		return &modelRepo.UserRedis{}
	}

	return &modelRepo.UserRedis{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: user.CreatedAt.UnixNano(),
		UpdatedAtNs: user.UpdatedAt.Time.UnixNano(),
	}
}

// ToUserFromRedisRepo converts user from redis repo to service user
func ToUserFromRedisRepo(userRedis *modelRepo.UserRedis) *model.User {
	if userRedis == nil {
		return &model.User{}
	}

	return &model.User{
		ID:        userRedis.ID,
		Name:      userRedis.Name,
		Email:     userRedis.Email,
		Role:      userRedis.Role,
		CreatedAt: time.Unix(0, userRedis.CreatedAtNs),
		UpdatedAt: sql.NullTime{
			Time:  time.Unix(0, userRedis.UpdatedAtNs),
			Valid: true,
		},
	}
}
