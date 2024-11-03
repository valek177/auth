package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/valek177/auth/grpc/pkg/user_v1"
	"github.com/valek177/auth/internal/api/auth"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *user_v1.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		time  = time.Now()
		role  = user_v1.Role_USER.String()

		serviceErr = fmt.Errorf("service error")

		req = &user_v1.GetUserRequest{
			Id: id,
		}

		user = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: time,
		}

		res = &user_v1.GetUserResponse{
			User: &user_v1.User{
				Id: id,
				UserInfo: &user_v1.UserInfo{
					Name:  wrapperspb.String(name),
					Email: wrapperspb.String(email),
					Role:  user_v1.Role_USER,
				},
				CreatedAt: timestamppb.New(time),
			},
		}
	)

	testsSuccessful := []struct {
		name            string
		args            args
		want            *user_v1.GetUserResponse
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
				mock.GetUserMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
		},
	}
	testsErrors := []struct {
		name            string
		args            args
		want            *user_v1.GetUserResponse
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
				mock.GetUserMock.Expect(ctx, id).Return(nil, serviceErr)
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
			err:  fmt.Errorf("unable to get user: empty request"),
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

			user, err := api.GetUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, user)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMock(mc)
			api := auth.NewImplementation(authServiceMock)

			_, err := api.GetUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
