package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	descUser "github.com/valek177/auth/internal/api/user"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *user_v1.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id         = gofakeit.Int64()
		name       = gofakeit.Name()
		email      = gofakeit.Email()
		password   = gofakeit.Password(true, false, false, false, false, 7)
		serviceErr = fmt.Errorf("service error")

		req = &user_v1.CreateUserRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            user_v1.Role_USER,
		}
		res = &user_v1.CreateUserResponse{
			Id: id,
		}
	)

	testsSuccessful := []struct {
		name            string
		args            args
		want            *user_v1.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Set(func(_ context.Context, _ *model.NewUser,
				) (i1 int64, err error) {
					return id, nil
				})
				return mock
			},
		},
	}
	testsErrors := []struct {
		name            string
		args            args
		want            *user_v1.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Set(func(_ context.Context, _ *model.NewUser) (i1 int64, err error) {
					return 0, serviceErr
				})
				return mock
			},
		},
		{
			name: "error: validation error (empty request)",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  fmt.Errorf("unable to create user: empty request"),
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range testsSuccessful {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := descUser.NewImplementation(userServiceMock)

			newID, err := api.CreateUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, newID)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := descUser.NewImplementation(userServiceMock)

			_, err := api.CreateUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
