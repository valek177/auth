package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/api/auth"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *user_v1.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, false, false, false, false, 7)

		serviceErr = fmt.Errorf("service error")

		req = &user_v1.CreateUserRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            user_v1.Role_USER,
		}

		newUser = &model.NewUser{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: password,
			Role:            user_v1.Role_USER.String(),
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
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case 1",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(id, nil)
				return mock
			},
		},
	}
	testsErrors := []struct {
		name            string
		args            args
		want            *user_v1.CreateUserResponse
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "service error case 1",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, newUser).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range testsSuccessful {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := auth.NewImplementation(authServiceMock)

			newID, err := api.CreateUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, newID)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := auth.NewImplementation(authServiceMock)

			_, err := api.CreateUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "service error")
		})
	}
}
