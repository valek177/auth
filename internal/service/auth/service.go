package auth

import (
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
	"github.com/valek177/platform-common/pkg/client/db"
)

type serv struct {
	authRepository  repository.AuthRepository
	logRepository   repository.LogRepository
	redisRepository repository.UserRedisRepository
	txManager       db.TxManager
}

// NewService creates new service with settings
func NewService(
	authRepository repository.AuthRepository,
	logRepository repository.LogRepository,
	redisRepository repository.UserRedisRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository:  authRepository,
		logRepository:   logRepository,
		redisRepository: redisRepository,
		txManager:       txManager,
	}
}
