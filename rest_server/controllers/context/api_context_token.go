package context

type TokenInfo struct {
	AccessToken           string `json:"access_token" validate:"required"`
	AccessTokenExpiredAt  string `json:"access_token_expired_at" validate:"required"`
	RefreshToken          string `json:"refresh_token" validate:"required"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at" validate:"required"`
}
