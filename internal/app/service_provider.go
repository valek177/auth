package app

import (
	"context"

	"github.com/IBM/sarama"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	accessImpl "github.com/valek177/auth/internal/api/access"
	authImpl "github.com/valek177/auth/internal/api/auth"
	userImpl "github.com/valek177/auth/internal/api/user"
	"github.com/valek177/auth/internal/client/kafka"
	kafkaConsumer "github.com/valek177/auth/internal/client/kafka/consumer"
	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/config/env"
	"github.com/valek177/auth/internal/repository"
	accessRepository "github.com/valek177/auth/internal/repository/access"
	logRepo "github.com/valek177/auth/internal/repository/log"
	redisRepo "github.com/valek177/auth/internal/repository/redis"
	userRepository "github.com/valek177/auth/internal/repository/user"
	"github.com/valek177/auth/internal/service"
	accessService "github.com/valek177/auth/internal/service/access"
	authService "github.com/valek177/auth/internal/service/auth"
	userSaverConsumer "github.com/valek177/auth/internal/service/consumer/user_saver"
	userService "github.com/valek177/auth/internal/service/user"
	"github.com/valek177/auth/internal/utils"
	cache "github.com/valek177/platform-common/pkg/client/cache"
	redisConfig "github.com/valek177/platform-common/pkg/client/cache/config"
	redis "github.com/valek177/platform-common/pkg/client/cache/redis"
	"github.com/valek177/platform-common/pkg/client/db"
	"github.com/valek177/platform-common/pkg/client/db/pg"
	"github.com/valek177/platform-common/pkg/client/db/transaction"
	"github.com/valek177/platform-common/pkg/closer"
)

type serviceProvider struct {
	pgConfig           config.PGConfig
	grpcConfig         config.GRPCConfig
	httpConfig         config.HTTPConfig
	redisConfig        redisConfig.RedisConfig
	swaggerConfig      config.SwaggerConfig
	tokenRefreshConfig config.TokenConfig
	tokenAccessConfig  config.TokenConfig
	prometheusConfig   config.PrometheusConfig

	kafkaConsumerConfig  config.KafkaConsumerConfig
	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler

	userSaverConsumer service.ConsumerService

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	tokenAccess  utils.Token
	tokenRefresh utils.Token

	userRepository   repository.UserRepository
	accessRepository repository.AccessRepository
	logRepository    repository.LogRepository
	redisRepository  repository.UserRedisRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *userImpl.Implementation
	authImpl   *authImpl.Implementation
	accessImpl *accessImpl.Implementation
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

// HTTPConfig returns HTTP config
func (s *serviceProvider) HTTPConfig() (config.HTTPConfig, error) {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			return nil, err
		}

		s.httpConfig = cfg
	}

	return s.httpConfig, nil
}

// SwaggerConfig return swagger config
func (s *serviceProvider) SwaggerConfig() (config.SwaggerConfig, error) {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			return nil, err
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig, nil
}

// TokenAccessConfig returns token access config
func (s *serviceProvider) TokenAccessConfig() (config.TokenConfig, error) {
	if s.tokenAccessConfig == nil {
		cfg, err := env.NewAccessTokenConfig()
		if err != nil {
			return nil, err
		}

		s.tokenAccessConfig = cfg
	}

	return s.tokenAccessConfig, nil
}

// TokenRefreshConfig returns token refresh config
func (s *serviceProvider) TokenRefreshConfig() (config.TokenConfig, error) {
	if s.tokenRefreshConfig == nil {
		cfg, err := env.NewRefreshTokenConfig()
		if err != nil {
			return nil, err
		}

		s.tokenRefreshConfig = cfg
	}

	return s.tokenRefreshConfig, nil
}

// RedisConfig returns redis config
func (s *serviceProvider) RedisConfig() (redisConfig.RedisConfig, error) {
	if s.redisConfig == nil {
		cfg, err := env.NewRedisConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig, nil
}

// KafkaConsumerConfig returns config for kafka consumer
func (s *serviceProvider) KafkaConsumerConfig() (config.KafkaConsumerConfig, error) {
	if s.kafkaConsumerConfig == nil {
		cfg, err := env.NewKafkaConsumerConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig, nil
}

// PrometheusConfig returns prometheus config
func (s *serviceProvider) PrometheusConfig() (config.PrometheusConfig, error) {
	if s.prometheusConfig == nil {
		cfg, err := env.NewPrometheusConfig()
		if err != nil {
			return nil, err
		}

		s.prometheusConfig = cfg
	}

	return s.prometheusConfig, nil
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

// RedisPool creates new redis pool
func (s *serviceProvider) RedisPool() (*redigo.Pool, error) {
	if s.redisPool == nil {
		redisConfig, err := s.RedisConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		s.redisPool = &redigo.Pool{
			MaxIdle:     redisConfig.MaxIdle(),
			IdleTimeout: redisConfig.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", redisConfig.Address())
			},
		}
	}

	return s.redisPool, nil
}

// RedisClient returns redis client
func (s *serviceProvider) RedisClient() (cache.RedisClient, error) {
	if s.redisClient == nil {
		pool, err := s.RedisPool()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		cfg, err := s.RedisConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		s.redisClient = redis.NewClient(pool, cfg)
	}

	return s.redisClient, nil
}

// TokenRefresh returns refresh token
func (s *serviceProvider) TokenRefresh() (utils.Token, error) {
	if s.tokenRefresh == nil {
		cfg, err := s.TokenRefreshConfig()
		if err != nil {
			return nil, err
		}
		s.tokenRefresh = utils.NewToken(cfg)
	}

	return s.tokenRefresh, nil
}

// TokenAccess returns access token
func (s *serviceProvider) TokenAccess() (utils.Token, error) {
	if s.tokenAccess == nil {
		cfg, err := s.TokenAccessConfig()
		if err != nil {
			return nil, err
		}
		s.tokenAccess = utils.NewToken(cfg)
	}

	return s.tokenAccess, nil
}

// UserRedisRepository returns redis repository
func (s *serviceProvider) UserRedisRepository() (
	repository.UserRedisRepository, error,
) {
	if s.redisRepository == nil {
		client, err := s.RedisClient()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		config, err := s.RedisConfig()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		s.redisRepository = redisRepo.NewUserRedisRepository(client, config)
	}

	return s.redisRepository, nil
}

// UserRepository returns new UserRepository
func (s *serviceProvider) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if s.userRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.userRepository = userRepository.NewRepository(dbClient)
	}

	return s.userRepository, nil
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

// AccessRepository returns access repository
func (s *serviceProvider) AccessRepository(ctx context.Context) (
	repository.AccessRepository, error,
) {
	if s.accessRepository == nil {
		dbClient, err := s.DBClient(ctx)
		if err != nil {
			return nil, err
		}
		s.accessRepository = accessRepository.NewRepository(dbClient)
	}

	return s.accessRepository, nil
}

// UserService returns new UserService
func (s *serviceProvider) UserService(ctx context.Context) (service.UserService, error) {
	if s.userService == nil {
		userRepo, err := s.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		logRepo, err := s.LogRepository(ctx)
		if err != nil {
			return nil, err
		}
		redisRepo, err := s.UserRedisRepository()
		if err != nil {
			return nil, err
		}
		txManager, err := s.TxManager(ctx)
		if err != nil {
			return nil, err
		}
		s.userService = userService.NewService(
			userRepo, logRepo, redisRepo, txManager,
		)
	}

	return s.userService, nil
}

// AuthService returns new AuthService
func (s *serviceProvider) AuthService(ctx context.Context) (service.AuthService, error) {
	if s.authService == nil {
		userRepo, err := s.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		tokenAccess, err := s.TokenAccess()
		if err != nil {
			return nil, err
		}
		tokenRefresh, err := s.TokenRefresh()
		if err != nil {
			return nil, err
		}
		s.authService = authService.NewService(
			userRepo,
			tokenRefresh,
			tokenAccess,
		)
	}

	return s.authService, nil
}

// AccessService returns new AccessService
func (s *serviceProvider) AccessService(ctx context.Context) (service.AccessService, error) {
	if s.accessService == nil {
		accessRepo, err := s.AccessRepository(ctx)
		if err != nil {
			return nil, err
		}
		tokenAccess, err := s.TokenAccess()
		if err != nil {
			return nil, err
		}
		s.accessService = accessService.NewService(accessRepo, tokenAccess)
	}

	return s.accessService, nil
}

// UserImpl returns new User Service implementation
func (s *serviceProvider) UserImpl(ctx context.Context) (*userImpl.Implementation, error) {
	if s.userImpl == nil {
		userServ, err := s.UserService(ctx)
		if err != nil {
			return nil, err
		}
		s.userImpl = userImpl.NewImplementation(userServ)
	}

	return s.userImpl, nil
}

// AuthImpl returns new Auth Service implementation
func (s *serviceProvider) AuthImpl(ctx context.Context) (*authImpl.Implementation, error) {
	if s.authImpl == nil {
		authServ, err := s.AuthService(ctx)
		if err != nil {
			return nil, err
		}
		s.authImpl = authImpl.NewImplementation(authServ)
	}

	return s.authImpl, nil
}

// AccessImpl returns new Access Service implementation
func (s *serviceProvider) AccessImpl(ctx context.Context) (*accessImpl.Implementation, error) {
	if s.accessImpl == nil {
		accessServ, err := s.AccessService(ctx)
		if err != nil {
			return nil, err
		}
		s.accessImpl = accessImpl.NewImplementation(accessServ)
	}

	return s.accessImpl, nil
}

// UserSaverConsumer returns user consumer service
func (s *serviceProvider) UserSaverConsumer(ctx context.Context) (service.ConsumerService, error) {
	if s.userSaverConsumer == nil {
		userRepo, err := s.UserRepository(ctx)
		if err != nil {
			return nil, err
		}
		consumer, err := s.Consumer()
		if err != nil {
			return nil, err
		}
		s.userSaverConsumer = userSaverConsumer.NewService(
			userRepo,
			consumer,
		)
	}

	return s.userSaverConsumer, nil
}

// Consumer returns kafka consumer
func (s *serviceProvider) Consumer() (kafka.Consumer, error) {
	if s.consumer == nil {
		group, err := s.ConsumerGroup()
		if err != nil {
			return nil, err
		}
		s.consumer = kafkaConsumer.NewConsumer(
			group,
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer, nil
}

// ConsumerGroup returns consumer group
func (s *serviceProvider) ConsumerGroup() (sarama.ConsumerGroup, error) {
	if s.consumerGroup == nil {
		cfg, err := s.KafkaConsumerConfig()
		if err != nil {
			return nil, err
		}
		consumerGroup, err := sarama.NewConsumerGroup(
			cfg.Brokers(),
			cfg.GroupID(),
			cfg.Config(),
		)
		if err != nil {
			return nil, err
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup, nil
}

// ConsumerGroupHandler returns consumer group handler
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}
