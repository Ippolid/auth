package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Ippolid/auth/internal/service/mocks"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	repoMocks "github.com/Ippolid/auth/internal/repository/mocks"
	"github.com/Ippolid/auth/internal/service/auth"
	"github.com/Ippolid/platform_libary/pkg/db" // Предполагаемый путь к мокам
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
	)
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		authRepositoryMock authRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "Delete" || log.Ctx != Ctxstring {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (*model.User, error) {
					return nil, fmt.Errorf("user not found")
				})
				return mock
			},

			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "DeleteUser error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			t.Cleanup(mc.Finish)

			authRepoMock := tt.authRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			service := auth.NewMockService(authRepoMock, txManagerMock)

			err := service.Delete(tt.args.ctx, tt.args.id)
			require.ErrorIs(t, err, tt.err)
		})
	}
}
