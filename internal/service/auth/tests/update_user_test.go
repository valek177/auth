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
	"github.com/valek177/auth/internal/service/auth"
	"github.com/valek177/platform-common/pkg/client/db"
	dbMocks "github.com/valek177/platform-common/pkg/client/db/mocks"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
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
		authRepositoryMock  authRepositoryMockFunc
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
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUser).Return(nil)
				mock.GetUserMock.Expect(ctx, id).Return(&model.User{}, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(ctx context.Context, record *model.Record,
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
				mock.CreateUserMock.Set(func(ctx context.Context, user *model.User) (err error) {
					return nil
				})
				mock.DeleteUserMock.Set(func(ctx context.Context, id int64) (err error) {
					return nil
				})

				mock.SetExpireUserMock.Set(func(ctx context.Context, id int64) (err error) {
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
				req: updateUser,
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUser).Return(repoErr)
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
			name: "error: validation error (empty model)",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  fmt.Errorf("unable to update user: empty model"),
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
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
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
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
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUser).Return(nil)
				mock.GetUserMock.Expect(ctx, id).Return(&model.User{}, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(ctx context.Context,
					record *model.Record,
				) (i1 int64, err error) {
					return 0, fmt.Errorf("create record on update error")
				})
				return mock
			},
			txManagerMock: txManagerFunc,
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.DeleteUserMock.Set(func(ctx context.Context, id int64) (err error) {
					return nil
				})
				mock.CreateUserMock.Set(func(ctx context.Context, user *model.User) (err error) {
					return nil
				})
				mock.SetExpireUserMock.Set(func(ctx context.Context, id int64) (err error) {
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

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
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

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			redisRepositoryMock := tt.redisRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			err := service.UpdateUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
