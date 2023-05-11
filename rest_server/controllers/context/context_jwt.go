package context

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

type TokenType int

const (
	AccessT TokenType = iota
	RefreshT
)

// Redis에 key로 사용될 텍스트 모음
var TokenTypeText = map[TokenType]string{
	AccessT:  "ACCESS",
	RefreshT: "REFRESH",
}

var UuidTypeText = map[TokenType]string{
	AccessT:  "access_uuid",
	RefreshT: "refresh_uuid",
}

// Jwt 토큰 정보
type JwtInfo struct {
	Idx int64 `json:"idx,omitempty"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpireDt   int64  `json:"access_token_expire_dt"`
	RtExpireDt   int64  `json:"refresh_token_expire_dt"`
}

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func GetTokenExpiryperiod(loginType LoginType) (int64, int64) {
	confAuth := config.GetInstance().Auth
	switch loginType {
	case AppLogin:
		return confAuth.AppAccessTokenExpiryPeriod * int64(time.Hour), confAuth.AppRefreshTokenExpiryPeriod * int64(time.Hour)
	case WebAccountLogin:
		return confAuth.WebAccessTokenExpiryPeriod * int64(time.Hour), confAuth.WebRefreshTokenExpiryPeriod * int64(time.Hour)
	case CustomerLogin:
		return confAuth.CustomerAccessTokenExpiryPeriod * int64(time.Hour), confAuth.CustomerRefreshTokenExpiryPeriod * int64(time.Hour)
	}
	return 0, 0
}
