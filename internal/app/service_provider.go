package app

import (
	"context"

	"github.com/valek177/auth/internal/api/auth"
	"github.com/valek177/auth/internal/client/db"
	"github.com/valek177/auth/internal/client/db/pg"
	"github.com/valek177/auth/internal/client/db/transaction"
	"github.com/valek177/auth/internal/closer"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/config/env"
	"github.com/valek177/auth/internal/repository"
	authRepository "github.com/valek177/auth/internal/repository/auth"
	logRepo "github.com/valek177/auth/internal/repository/log"
	"github.com/valek177/auth/internal/service"
	authService "github.com/valek177/auth/internal/service/auth"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	authRepository repository.AuthRepository
	logRepository  repository.LogRepository

	authService service.AuthService

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig returns new PGConfig
func (s *serviceProvider) PGConfig() (config.PGConfig, error) {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			return nil, err
		}

		s.pgConfig = cfg
	}

	return s.pgConfig, nil
}

// GRPCConfig returns new GRPCConfig
func (s *serviceProvider) GRPCConfig() (config.GRPCConfig, error) {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			return nil, err
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig, nil
}

// DBClient returns new db client
func (s *serviceProvider) DBClient(ctx context.Context) (db.Client, error) {
	if s.dbClient == nil {
		pgConfig, err := s.PGConfig()
		if err != nil {
			return nil, err
		}
		cl, err := pg.New(ctx, pgConfig.DSN())
		if err != nil {
			return nil, err
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			return nil, err
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient, nil
}

// TxManager returns new db TxManager
func (s *serviceProvider) TxManager(ctx context.Context) (db.TxManager, error) {
	if s.txManager == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.txManager = transaction.NewTransactionManager(dbClient.DB())
	}

	return s.txManager, nil
}

// AuthRepository returns new AuthRepository
func (s *serviceProvider) AuthRepository(ctx context.Context) (repository.AuthRepository, error) {
	if s.authRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.authRepository = authRepository.NewRepository(dbClient)
	}

	return s.authRepository, nil
}

// LogRepository returns new LogRepository
func (s *serviceProvider) LogRepository(ctx context.Context) (repository.LogRepository, error) {
	if s.logRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.logRepository = logRepo.NewRepository(dbClient)
	}

	return s.logRepository, nil
}

// AuthService returns new AuthService
func (s *serviceProvider) AuthService(ctx context.Context) (service.AuthService, error) {
	if s.authService == nil {
		authRepo, err := s.AuthRepository(ctx)
		if err != nil {
			return nil, err
		}
		logRepo, err := s.LogRepository(ctx)
		if err != nil {
			return nil, err
		}
		txManager, err := s.TxManager(ctx)
		if err != nil {
			return nil, err
		}
		s.authService = authService.NewService(
			authRepo, logRepo, txManager,
		)
	}

	return s.authService, nil
}

// AuthImpl returns new Auth Service implementation
func (s *serviceProvider) AuthImpl(ctx context.Context) (*auth.Implementation, error) {
	if s.authImpl == nil {
		authServ, err := s.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		s.authImpl = auth.NewImplementation(authServ)
	}

	return s.authImpl, nil
}
