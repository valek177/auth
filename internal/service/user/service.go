package user

import (
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
	"github.com/valek177/platform-common/pkg/client/db"
)

type serv struct {
	userRepository  repository.UserRepository
	logRepository   repository.LogRepository
	redisRepository repository.UserRedisRepository
	txManager       db.TxManager
}

// NewService creates new service with settings
func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	redisRepository repository.UserRedisRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository:  userRepository,
		logRepository:   logRepository,
		redisRepository: redisRepository,
		txManager:       txManager,
	}
}
