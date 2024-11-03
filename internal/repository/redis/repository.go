package pg

import (
	"context"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/repository/redis/converter"
	modelRepo "github.com/valek177/auth/internal/repository/redis/model"
	"github.com/valek177/platform-common/pkg/client/cache"
	"github.com/valek177/platform-common/pkg/client/cache/config"
)

type repo struct {
	cl     cache.RedisClient
	config config.RedisConfig
}

// NewUserRedisRepository returns new user redis repository
func NewUserRedisRepository(cl cache.RedisClient, config config.RedisConfig,
) repository.UserRedisRepository {
	return &repo{cl: cl, config: config}
}

// CreateUser creates user record in Redis
func (r *repo) CreateUser(ctx context.Context, user *model.User) error {
	userRedis := converter.ToRedisRepoFromUser(user)

	idStr := strconv.FormatInt(userRedis.ID, 10)
	err := r.cl.Ping(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = r.cl.HashSet(ctx, idStr, userRedis)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetUser returns user from redis
func (r *repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.cl.HGetAll(ctx, idStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(values) == 0 {
		return nil, model.ErrorUserNotFound
	}

	var userRedis modelRepo.UserRedis
	err = redigo.ScanStruct(values, &userRedis)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return converter.ToUserFromRedisRepo(&userRedis), nil
}

// DeleteUser deletes user in redis
func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)

	err := r.cl.Delete(ctx, idStr)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// SetExpireUser sets expire time of user
func (r *repo) SetExpireUser(ctx context.Context, id int64) error {
	idStr := strconv.FormatInt(id, 10)

	err := r.cl.Expire(ctx, idStr, r.config.ElementTTL())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
