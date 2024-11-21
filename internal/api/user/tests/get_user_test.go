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
	descUser "github.com/valek177/auth/internal/api/user"
	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/service"
	serviceMocks "github.com/valek177/auth/internal/service/mocks"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

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

			user, err := api.GetUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, user)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := descUser.NewImplementation(userServiceMock)

			_, err := api.GetUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
