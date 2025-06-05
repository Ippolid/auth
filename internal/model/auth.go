package model

// LoginRequest структура запроса для входа в систему
type LoginRequest struct {
	Username string
	Password string
}

// LoginResponse структура ответа при входе в систему
type LoginResponse struct {
	Token string
}

// GetRefreshTokenRequest структура запроса для получения нового refresh-токена
type GetRefreshTokenRequest struct {
	OldToken string
}

// GetRefreshTokenResponse структура ответа при получении нового refresh-токена
type GetRefreshTokenResponse struct {
	RefreshToken string
}

// GetAccessTokenRequest структура запроса для получения access-токена по refresh-токену
type GetAccessTokenRequest struct {
	RefreshToken string
}

// GetAccessTokenResponse структура ответа при получении access-токена
type GetAccessTokenResponse struct {
	AccessToken string
}

// CheckRequest структура запроса для проверки доступности конечной точки
type CheckRequest struct {
	EndpointAddress string
}
