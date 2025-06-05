package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//	func (s *serv) GetAccessToken(ctx context.Context, req model.GetAccessTokenRequest) (*model.GetAccessTokenResponse, error) {
//		claims, err := utils.VerifyToken(req.RefreshToken, []byte(s.token.RefreshToken()))
//		if err != nil {
//			return nil, status.Errorf(codes.Aborted, "invalid refresh token")
//		}
//		fmt.Println(claims)
//		role, errCache := s.cache.GetRole(ctx, claims.Username)
//		fmt.Println(role, errCache)
//		if errCache != nil || role == nil {
//			if errors.Is(errCache, model.ErrUserNotFound) {
//				err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
//
//					var errTx error
//					*role, errTx = s.authRepository.GetUserRole(ctx, claims.Username)
//					if errTx != nil {
//						return fmt.Errorf("error getting user profile: %w", errTx)
//					}
//					fmt.Println(*role)
//
//					errTx = s.authRepository.MakeLog(ctx, model.Log{
//						Method:    "Get user Role",
//						CreatedAt: time.Now(),
//						Ctx:       fmt.Sprintf("%v", ctx),
//					})
//					if errTx != nil {
//						return fmt.Errorf("error creating log: %w", errTx)
//					}
//
//					if errTx = s.cache.CreateRole(ctx, claims.Username, *role); errTx != nil {
//						return fmt.Errorf("error caching user profile: %w", errTx)
//					}
//
//					return nil
//				})
//				if err != nil {
//					return nil, err
//				}
//			}
//			// если ошибка другая — возвращаем ошибку кеша
//			return nil, fmt.Errorf("error with cache: %w", errCache)
//		}
//
//		accessToken, err := utils.GenerateToken(model.UserInfoJwt{
//			Username: claims.Username,
//			// Это пример, в реальности роль должна браться из базы или кэша
//			Role: *role,
//		},
//			[]byte(s.token.AccessToken()),
//			accessTokenExpiration,
//		)
//		if err != nil {
//			return nil, err
//		}
//
//		return &model.GetAccessTokenResponse{AccessToken: accessToken}, nil
//	}
func (s *serv) GetAccessToken(ctx context.Context, req model.GetAccessTokenRequest) (*model.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.RefreshToken, []byte(s.token.RefreshToken()))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	var userRole bool
	rolePtr, errCache := s.cache.GetRole(ctx, claims.Username)

	if errCache != nil {
		if errors.Is(errCache, model.ErrUserNotFound) {
			err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
				var errTx error
				userRole, errTx = s.authRepository.GetUserRole(ctx, claims.Username)
				if errTx != nil {
					return fmt.Errorf("error getting user role: %w", errTx)
				}

				errTx = s.authRepository.MakeLog(ctx, model.Log{
					Method:    "Get user Role",
					CreatedAt: time.Now(),
					Ctx:       fmt.Sprintf("%v", ctx),
				})
				if errTx != nil {
					return fmt.Errorf("error creating log: %w", errTx)
				}

				if errTx = s.cache.CreateRole(ctx, claims.Username, userRole); errTx != nil {
					return fmt.Errorf("error caching user role: %w", errTx)
				}

				return nil
			})
			if err != nil {
				return nil, err
			}
		} else {
			// Только если ошибка не ErrUserNotFound, возвращаем ошибку кеша
			return nil, fmt.Errorf("error with cache: %w", errCache)
		}
	} else if rolePtr != nil {
		// Если роль успешно получена из кеша
		userRole = *rolePtr
	} else {
		return nil, fmt.Errorf("роль пользователя не найдена")
	}

	accessToken, err := utils.GenerateToken(model.UserInfoJwt{
		Username: claims.Username,
		Role:     userRole,
	},
		[]byte(s.token.AccessToken()),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &model.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
