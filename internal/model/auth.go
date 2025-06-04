package model

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
}

type GetRefreshTokenRequest struct {
	OldToken string
}

type GetRefreshTokenResponse struct {
	RefreshToken string
}

type GetAccessTokenRequest struct {
	RefreshToken string
}

type GetAccessTokenResponse struct {
	AccessToken string
}

type CheckRequest struct {
	EndpointAddress string
}
