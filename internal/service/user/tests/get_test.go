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
	"github.com/Ippolid/auth/internal/service/user"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

var Ctxstring = "context.Background"

func TestGet(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager
	type cacheMockFunc func(mc *minimock.Controller) repository.CacheInterface

	type args struct {
		ctx context.Context
		id  int64
	}

	ctx := context.Background()
	id := gofakeit.Int64()
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, 10)
	role := gofakeit.Bool()
	createdAt := time.Now()
	repoErr := fmt.Errorf("repo error")

	expectedUser := &model.User{
		ID: id,
		User: model.UserInfo{
			Name:  &name,
			Email: &email,
		},
		Role:      role,
		Password:  password,
		CreatedAt: createdAt,
	}

	tests := []struct {
		name               string
		args               args
		wantUser           *model.User
		wantErr            error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
		cacheMock          cacheMockFunc
	}{
		{
			name: "user found in cache",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantUser: expectedUser,
			wantErr:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				// Не должен вызываться
				return repoMocks.NewUserRepositoryMock(mc)
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				// Не должен вызываться
				return mocks.NewTxManagerMock(mc)
			},
			cacheMock: func(mc *minimock.Controller) repository.CacheInterface {
				mock := repoMocks.NewCacheInterfaceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(expectedUser, nil)
				return mock
			},
		},
		{
			name: "user not in cache, repo returns error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantUser: nil,
			wantErr:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, repoErr)
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
				mock.GetMock.Expect(ctx, id).Return(nil, model.ErrUserNotFound)
				return mock
			},
		},
		{
			name: "cache returns unexpected error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			wantUser: nil,
			wantErr:  fmt.Errorf("error with cache: %w", repoErr),
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				// Не должен вызываться
				return repoMocks.NewUserRepositoryMock(mc)
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				// Не должен вызываться
				return mocks.NewTxManagerMock(mc)
			},
			cacheMock: func(mc *minimock.Controller) repository.CacheInterface {
				mock := repoMocks.NewCacheInterfaceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
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
			service := user.NewService(userRepoMock, txManagerMock, cacheMock)

			user, err := service.Get(tt.args.ctx, tt.args.id)
			if tt.wantErr != nil {
				require.Error(t, err)
				if tt.name == "cache returns unexpected error" {
					require.Contains(t, err.Error(), "error with cache")
				} else {
					require.ErrorIs(t, err, tt.wantErr)
				}
				require.Equal(t, tt.wantUser, user)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUser, user)
			}
		})
	}
}
