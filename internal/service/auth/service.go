package auth

import (
	"github.com/valek177/auth/internal/client/db"
	"github.com/valek177/auth/internal/repository"
	"github.com/valek177/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(
	authRepository repository.AuthRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
