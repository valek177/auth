package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/api/auth"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *user_v1.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		role = user_v1.Role_USER.String()

		serviceErr = fmt.Errorf("service error")

		req = &user_v1.UpdateUserRequest{
			Id:   id,
			Name: wrapperspb.String(name),
			Role: user_v1.Role_USER,
		}

		updateUser = &model.UpdateUserInfo{
			ID:   id,
			Name: &name,
			Role: &role,
		}

		res = &emptypb.Empty{}
	)

	testsSuccessful := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUser).Return(nil)
				return mock
			},
		},
	}
	testsErrors := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, updateUser).Return(serviceErr)
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
			err:  fmt.Errorf("unable to update user: empty request"),
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
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

			res, err := api.UpdateUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, res)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := auth.NewImplementation(authServiceMock)

			_, err := api.UpdateUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
