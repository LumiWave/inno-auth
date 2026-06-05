package auth

import (
	"errors"

	"github.com/LumiWave/baseutil/log"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

const (
	AppleJWKsEndpoint = "https://appleid.apple.com/auth/keys"
)

type AppleUser struct {
	EA     string `json:"email"`
	UserID string `json:"sub"`
}

type OauthApple struct {
	SocialType int64
}

func NewOauthApple() *OauthApple {
	return new(OauthApple)
}

func (o *OauthApple) GetSocialType() int64 {
	return o.SocialType
}

// Apple 은 identity token(JWT) 자체를 Apple 공개키(JWKs)로 검증한다.
func (o *OauthApple) VerifySocialKey(socialKey string) (string, string, error) {
	appleUser := new(AppleUser)
	if err := getAppleJWKsVerify(socialKey, appleUser); err != nil {
		return "", "", err
	}
	return appleUser.UserID, appleUser.EA, nil
}

func getAppleJWKsVerify(socialKey string, appleUser *AppleUser) error {
	jwks, err := keyfunc.Get(AppleJWKsEndpoint, keyfunc.Options{})
	if err != nil {
		log.Errorf("apple keyfunc.Get err : %v", err)
		return err
	}

	// 토큰 파싱 및 서명 검증
	token, err := jwt.Parse(socialKey, jwks.Keyfunc)
	if err != nil || !token.Valid {
		log.Errorf("apple not invalid token err : %v", err)
		if err == nil {
			err = errors.New("apple invalid token")
		}
		return err
	}

	// 클레임 추출
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("apple token.Claims fail")
		return errors.New("apple token claims parse fail")
	}

	// sub : Apple 고유 사용자 식별자
	if val, ok := claims["sub"].(string); ok {
		appleUser.UserID = val
	}

	// email : 최초 로그인 시에만 포함될 수 있음 (이후 로그인에는 없을 수 있음)
	if val, ok := claims["email"].(string); ok {
		appleUser.EA = val
	}

	return nil
}
