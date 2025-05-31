package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	repoMocks "github.com/Ippolid/auth/internal/repository/mocks"
	"github.com/Ippolid/auth/internal/service/mocks"
	user1 "github.com/Ippolid/auth/internal/service/user"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) repository.CacheInterface

	type args struct {
		ctx  context.Context
		user *model.User
	}

	var (
		ctx     = context.Background()
		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
		user    = &model.User{
			ID: id,
			User: model.UserInfo{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
			},
			Role:      gofakeit.Bool(),
			Password:  gofakeit.Password(true, true, true, true, false, 10),
			CreatedAt: time.Now(),
		}
	)

	tests := []struct {
		name               string
		args               args
		wantID             int64
		wantErr            error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
		cacheMock          cacheMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				user: user,
			},
			wantID:  id,
			wantErr: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, *user).Return(id, nil)
				mock.MakeLogMock.Set(func(_ context.Context, log model.Log) error {
					if log.Method != "Create" || log.Ctx != Ctxstring {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					return nil
				})
				mock.GetUserMock.Set(func(_ context.Context, _ int64) (*model.User, error) {
					return user, nil
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
			cacheMock: func(mc *minimock.Controller) repository.CacheInterface {
				mock := repoMocks.NewCacheInterfaceMock(mc)
				mock.CreateMock.Expect(ctx, id, *user).Return(nil)
				return mock
			},
		},
		{
			name: "CreateUser error case",
			args: args{
				ctx:  ctx,
				user: user,
			},
			wantID:  0,
			wantErr: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, *user).Return(int64(0), repoErr)
				// MakeLog и GetUser не должны вызываться при ошибке CreateUser
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
			cacheMock: func(mc *minimock.Controller) repository.CacheInterface {
				mock := repoMocks.NewCacheInterfaceMock(mc)

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			t.Cleanup(mc.Finish)

			userRepoMock := tt.userRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			cacheMock := tt.cacheMock(mc)
			service := user1.NewService(userRepoMock, txManagerMock, cacheMock)

			gotID, err := service.Create(tt.args.ctx, tt.args.user)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantID, gotID)
		})
	}
}
