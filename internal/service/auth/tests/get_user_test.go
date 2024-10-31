package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	repoMocks "github.com/valek177/auth/internal/repository/mocks"
	"github.com/valek177/auth/internal/service/auth"
	"github.com/valek177/platform-common/pkg/client/db"
	dbMocks "github.com/valek177/platform-common/pkg/client/db/mocks"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type redisRepositoryMockFunc func(mc *minimock.Controller) repository.UserRedisRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")

		user = &model.User{
			ID:        id,
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      user_v1.Role_USER.String(),
			CreatedAt: time.Now(),
		}
	)

	testsSuccessful := []struct {
		name                string
		args                args
		want                *model.User
		err                 error
		authRepositoryMock  authRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "success case (get user from database)",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: user,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.GetUserMock.Set(func(ctx context.Context, id int64) (
					up1 *model.User, err error,
				) {
					return nil, errors.New("no user in redis")
				})
				mock.SetExpireUserMock.Set(func(ctx context.Context, id int64) (
					err error,
				) {
					return nil
				})
				mock.CreateUserMock.Set(func(ctx context.Context, user *model.User,
				) (err error) {
					return nil
				})
				return mock
			},
		},
		{
			name: "success case (get user from redis)",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: user,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.GetUserMock.Set(func(ctx context.Context, id int64) (
					up1 *model.User, err error,
				) {
					return user, nil
				})
				mock.SetExpireUserMock.Set(func(ctx context.Context, id int64) (
					err error,
				) {
					return nil
				})
				return mock
			},
		},
	}
	testsErrors := []struct {
		name                string
		args                args
		want                *emptypb.Empty
		err                 error
		authRepositoryMock  authRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "repo error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.GetUserMock.Set(func(ctx context.Context, id int64) (
					up1 *model.User, err error,
				) {
					return nil, errors.New("no user in redis")
				})
				return mock
			},
		},
	}

	for _, tt := range testsSuccessful {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			res, err := service.GetUser(tt.args.ctx, tt.args.id)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			_, err := service.GetUser(tt.args.ctx, tt.args.id)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "repo error")
		})
	}
}
