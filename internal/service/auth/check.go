package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (s *serv) Check(ctx context.Context, request model.CheckRequest) error {
	// Извлекаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	// Проверяем наличие заголовка авторизации
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	// Проверяем формат заголовка авторизации
	if !strings.HasPrefix(authHeader[0], "Bearer ") {
		return errors.New("invalid authorization header format")
	}

	// Извлекаем токен доступа
	accessToken := strings.TrimPrefix(authHeader[0], "Bearer ")

	// Проверяем токен и извлекаем claims
	claims, err := utils.VerifyToken(accessToken, []byte(s.token.AccessToken()))
	if err != nil {
		return fmt.Errorf("access token is invalid: %w", err)
	}

	// Получаем карту доступных ролей для эндпоинтов
	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get accessible roles: %w", err)
	}

	fmt.Println(accessibleMap, err)

	// Проверяем, есть ли запрашиваемый эндпоинт в карте доступа
	requiredRole, ok := accessibleMap[request.EndpointAddress]
	if !ok {
		// Эндпоинт не требует проверки доступа
		return fmt.Errorf("endpoint address %s is not accessible", request.EndpointAddress)
	}

	// Проверяем, соответствует ли роль пользователя требуемой роли
	if requiredRole == claims.Role || claims.Role == "admin" {
		// Пользователь имеет нужную роль или является администратором
		return nil
	}

	// В доступе отказано
	return errors.New("access denied for this endpoint")
}
