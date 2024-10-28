package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/api/auth"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *user_v1.DeleteUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &user_v1.DeleteUserRequest{
			Id: id,
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
			name: "success case 1",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
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
			name: "service error case 1",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(serviceErr)
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

			res, err := api.DeleteUser(tt.args.ctx, tt.args.req)

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

			_, err := api.DeleteUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, "service error")
		})
	}
}
