package tests

import (
	"context"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	repoMocks "github.com/Ippolid/auth/internal/repository/mocks"
	"github.com/Ippolid/auth/internal/service/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		password  = gofakeit.Password(true, true, true, true, false, 10)
		role      = gofakeit.Bool()
		createdAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		expectedUser = &model.User{
			ID: id,
			User: model.UserInfo{
				Name:  name,
				Email: email,
			},
			Role:      role,
			Password:  password,
			CreatedAt: createdAt,
		}
	)
	// Важно: mc.Finish() должен быть вызван после завершения всех параллельных тестов
	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: expectedUser,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(expectedUser, nil)

				// Используем функцию сравнения вместо точного значения для model.Log
				mock.MakeLogMock.Set(func(ctx context.Context, log model.Log) (err error) {
					// Проверяем только интересующие нас поля
					if log.Method != "GET" || log.Ctx != "context.Background" {
						return fmt.Errorf("unexpected log entry: %+v", log)
					}
					// Не проверяем CreatedAt, так как оно динамическое
					return nil
				})

				return mock
			},
		},
		{
			name: "GetUser error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, id).Return(nil, repoErr)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mc := minimock.NewController(t)
			t.Cleanup(mc.Finish)

			authRepoMock := tt.authRepositoryMock(mc)
			service := auth.NewMockService(authRepoMock) // Используем конструктор сервиса

			user, err := service.Get(tt.args.ctx, tt.args.id)

			require.ErrorIs(t, err, tt.err) // Проверяем тип ошибки
			require.Equal(t, tt.want, user) // Сравниваем ожидаемый и фактический результат
		})
	}
}
