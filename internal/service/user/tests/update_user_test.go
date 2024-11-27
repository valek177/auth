package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	repoMocks "github.com/valek177/auth/internal/repository/mocks"
	descUser "github.com/valek177/auth/internal/service/user"
	"github.com/valek177/platform-common/pkg/client/db"
	dbMocks "github.com/valek177/platform-common/pkg/client/db/mocks"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type redisRepositoryMockFunc func(mc *minimock.Controller) repository.UserRedisRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UpdateUserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		role = user_v1.Role_USER.String()

		repoErr = errors.New("repo error")

		updateUser = &model.UpdateUserInfo{
			ID:   id,
			Name: &name,
			Role: &role,
		}

		res = &emptypb.Empty{}
	)

	txManagerFunc := func(mc *minimock.Controller) db.TxManager {
		mock := dbMocks.NewTxManagerMock(mc)
		mock.ReadCommittedMock.
			Set(func(ctx context.Context, f db.Handler) error { return f(ctx) })
		return mock
	}

	testsSuccessful := []struct {
		name                string
		args                args
		want                *emptypb.Empty
		err                 error
		userRepositoryMock  userRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: updateUser,
			},
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Set(func(_ context.Context, _ *model.UpdateUserInfo) (
					err error,
				) {
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (
					up1 *model.User, err error,
				) {
					return &model.User{}, nil
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(_ context.Context, _ *model.Record,
				) (int64, error) {
					return 0, nil
				})
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.
					Set(func(ctx context.Context, f db.Handler) error { return f(ctx) })
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.CreateUserMock.Set(func(_ context.Context, _ *model.User) (err error) {
					return nil
				})
				mock.DeleteUserMock.Set(func(_ context.Context, _ int64) (err error) {
					return nil
				})

				mock.SetExpireUserMock.Set(func(_ context.Context, _ int64) (err error) {
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
		userRepositoryMock  userRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "repo error",
			args: args{
				ctx: ctx,
				req: updateUser,
			},
			want: nil,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Set(func(_ context.Context, _ *model.UpdateUserInfo) (
					err error,
				) {
					return repoErr
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "repo error: auth get user error",
			args: args{
				ctx: ctx,
				req: updateUser,
			},
			want: nil,
			err:  fmt.Errorf("auth update user error"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Set(func(_ context.Context, _ *model.UpdateUserInfo) (
					err error,
				) {
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (
					up1 *model.User, err error,
				) {
					return &model.User{}, errors.New("auth update user error")
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "repo redis error: delete user error",
			args: args{
				ctx: ctx,
				req: updateUser,
			},
			want: nil,
			err:  fmt.Errorf("redis delete user error"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Set(func(_ context.Context, _ *model.UpdateUserInfo) (
					err error,
				) {
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (
					up1 *model.User, err error,
				) {
					return &model.User{ID: id}, nil
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.DeleteUserMock.Set(func(_ context.Context, _ int64) (err error) {
					return errors.New("redis delete user error")
				})
				return mock
			},
		},
		{
			name: "error: validation error (empty model)",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  fmt.Errorf("unable to update user: empty model"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "error: validation error (nothing to update)",
			args: args{
				ctx: ctx,
				req: &model.UpdateUserInfo{
					ID:   id,
					Name: nil,
					Role: nil,
				},
			},
			want: nil,
			err:  fmt.Errorf("unable to update user: nothing to update"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "repo error: create record",
			args: args{
				ctx: ctx,
				req: updateUser,
			},
			want: nil,
			err:  fmt.Errorf("create record on update error"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Set(func(_ context.Context, _ *model.UpdateUserInfo) (
					err error,
				) {
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (
					up1 *model.User, err error,
				) {
					return &model.User{}, nil
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(_ context.Context,
					_ *model.Record,
				) (i1 int64, err error) {
					return 0, fmt.Errorf("create record on update error")
				})
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.DeleteUserMock.Set(func(_ context.Context, _ int64) (err error) {
					return nil
				})
				mock.CreateUserMock.Set(func(_ context.Context, _ *model.User) (err error) {
					return nil
				})
				mock.SetExpireUserMock.Set(func(_ context.Context, _ int64) (err error) {
					return nil
				})
				return mock
			},
		},
	}

	for _, tt := range testsSuccessful {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := descUser.NewService(userRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			err := service.UpdateUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := descUser.NewService(userRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			err := service.UpdateUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
