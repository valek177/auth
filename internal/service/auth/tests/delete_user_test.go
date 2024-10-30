package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/valek177/auth/internal/model"
	"github.com/valek177/auth/internal/repository"
	repoMocks "github.com/valek177/auth/internal/repository/mocks"
	"github.com/valek177/auth/internal/service/auth"
	"github.com/valek177/platform-common/pkg/client/db"
	dbMocks "github.com/valek177/platform-common/pkg/client/db/mocks"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")

		res = &emptypb.Empty{}
	)

	txManagerFunc := func(mc *minimock.Controller) db.TxManager {
		mock := dbMocks.NewTxManagerMock(mc)
		mock.ReadCommittedMock.
			Set(func(ctx context.Context, f db.Handler) error { return f(ctx) })
		return mock
	}

	testsSuccessful := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		authRepositoryMock authRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,

			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
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
			txManagerMock: txManagerFunc,
		},
	}
	testsErrors := []struct {
		name               string
		args               args
		want               *emptypb.Empty
		err                error
		authRepositoryMock authRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "repo error",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			txManagerMock: txManagerFunc,
		},
		{
			name: "repo error: create record",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  fmt.Errorf("create record on delete error"),
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.CreateRecordMock.Set(func(ctx context.Context,
					record *model.Record,
				) (i1 int64, err error) {
					return 0, fmt.Errorf("create record on delete error")
				})
				return mock
			},
			txManagerMock: txManagerFunc,
		},
	}

	for _, tt := range testsSuccessful {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
				txManagerMock)

			err := service.DeleteUser(tt.args.ctx, tt.args.req)

			assert.Nil(t, err)
		})
	}

	for _, tt := range testsErrors {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authRepositoryMock := tt.authRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(authRepositoryMock, logRepositoryMock,
				txManagerMock)

			err := service.DeleteUser(tt.args.ctx, tt.args.req)

			assert.NotNil(t, err)
			assert.ErrorContains(t, err, tt.err.Error())
		})
	}
}
