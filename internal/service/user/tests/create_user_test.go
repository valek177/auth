package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/model"
	passLib "github.com/valek177/auth/internal/password"
	"github.com/valek177/auth/internal/repository"
	repoMocks "github.com/valek177/auth/internal/repository/mocks"
	"github.com/valek177/auth/internal/service/user"
	"github.com/valek177/platform-common/pkg/client/db"
	dbMocks "github.com/valek177/platform-common/pkg/client/db/mocks"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type redisRepositoryMockFunc func(mc *minimock.Controller) repository.UserRedisRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx     context.Context
		newUser *model.NewUser
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = int64(123)
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, false, false, false, 7)

		repoErr = fmt.Errorf("repo error")

		hash, _ = passLib.HashPassword(password)
		newUser = &model.NewUser{
			Name:            name,
			Email:           email,
			Password:        hash,
			PasswordConfirm: password,
			Role:            user_v1.Role_USER.String(),
		}
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
		want                int64
		err                 error
		userRepositoryMock  userRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				newUser: newUser,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(id, nil)
				mock.GetUserMock.Expect(ctx, id).Return(&model.User{ID: 123}, nil)
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
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, &model.User{ID: id, Name: ""}).Return(nil)
				mock.SetExpireUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
	}
	testsErrors := []struct {
		name                string
		args                args
		want                int64
		err                 error
		userRepositoryMock  userRepositoryMockFunc
		logRepositoryMock   logRepositoryMockFunc
		redisRepositoryMock redisRepositoryMockFunc
		txManagerMock       txManagerMockFunc
	}{
		{
			name: "repo error",
			args: args{
				ctx:     ctx,
				newUser: newUser,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(0, repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "error: validation error (empty model)",
			args: args{
				ctx:     ctx,
				newUser: nil,
			},
			want: 0,
			err:  fmt.Errorf("unable to create user: empty model"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "error: validation error (empty name)",
			args: args{
				ctx: ctx,
				newUser: &model.NewUser{
					Name:            "",
					Email:           "test",
					Password:        "test",
					PasswordConfirm: "test",
					Role:            user_v1.Role_USER.String(),
				},
			},
			want: 0,
			err:  fmt.Errorf("unable to create user: name is required"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "error: validation error (empty password)",
			args: args{
				ctx: ctx,
				newUser: &model.NewUser{
					Name:            "test",
					Email:           "test",
					Password:        "",
					PasswordConfirm: "test",
					Role:            user_v1.Role_USER.String(),
				},
			},
			want: 0,
			err:  fmt.Errorf("unable to create user: password is required"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "error: validation error (passwords do not match)",
			args: args{
				ctx: ctx,
				newUser: &model.NewUser{
					Name:            "test",
					Email:           "test",
					Password:        "123",
					PasswordConfirm: "test",
					Role:            user_v1.Role_USER.String(),
				},
			},
			want: 0,
			err:  fmt.Errorf("unable to create user: the passwords do not match"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "repo error: get user",
			args: args{
				ctx:     ctx,
				newUser: newUser,
			},
			want: 0,
			err:  fmt.Errorf("get user error"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(id, nil)
				mock.GetUserMock.Expect(ctx, id).Return(nil, fmt.Errorf("get user error"))
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "repo error: create record",
			args: args{
				ctx:     ctx,
				newUser: newUser,
			},
			want: 0,
			err:  fmt.Errorf("create record on create error"),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(id, nil)
				mock.GetUserMock.Expect(ctx, id).Return(&model.User{}, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(_ context.Context,
					_ *model.Record,
				) (i1 int64, err error) {
					return 0, fmt.Errorf("create record on create error")
				})
				return mock
			},
			redisRepositoryMock: func(mc *minimock.Controller) repository.UserRedisRepository {
				mock := repoMocks.NewUserRedisRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, &model.User{}).Return(nil)
				mock.SetExpireUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
			txManagerMock: txManagerFunc,
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

			service := user.NewService(userRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			newID, err := service.CreateUser(tt.args.ctx, tt.args.newUser)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, newID)
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

			service := user.NewService(userRepositoryMock, logRepositoryMock,
				redisRepositoryMock, txManagerMock)

			_, err := service.CreateUser(tt.args.ctx, tt.args.newUser)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
