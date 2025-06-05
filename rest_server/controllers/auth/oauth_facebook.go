package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/LumiWave/baseutil/log"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

const (
	FacebookUserInfoAPIEndpoint = "https://graph.facebook.com/me?fields=email&access_token="
	FacebookJWKsEndpoint        = "https://www.facebook.com/.well-known/oauth/openid/jwks/"
)

type FacebookUser struct {
	EA     string `json:"email"`
	UserID string `json:"id"`
}
type OauthFacebook struct {
	SocialType int64
}

func NewOauthFacebook() *OauthFacebook {
	return new(OauthFacebook)
}

func (o *OauthFacebook) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthFacebook) VerifySocialKey(socialKey string) (string, string, error) {
	userInfoResp, err := http.Get(FacebookUserInfoAPIEndpoint + url.QueryEscape(socialKey))
	if err != nil {
		return "", "", err
	}
	defer userInfoResp.Body.Close()

	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		return "", "", err
	}

	facebookUser := new(FacebookUser)
	json.Unmarshal(userInfo, &facebookUser)

	// ios sdk 17 으로 로그인 시는 에러나기 때문에 대응 코드 추가 한다.
	if len(facebookUser.UserID) == 0 {
		getJWKsVerify(socialKey, facebookUser)
	}
	return facebookUser.UserID, facebookUser.EA, nil
}

func getJWKsVerify(socialKey string, facebookUser *FacebookUser) {
	jwks, err := keyfunc.Get(FacebookJWKsEndpoint, keyfunc.Options{})
	if err != nil {
		log.Errorf("keyfunc.Get err : %v", err)
		return
	}

	// 토큰 파싱 및 검증
	token, err := jwt.Parse(socialKey, jwks.Keyfunc)
	if err != nil || !token.Valid {
		log.Errorf("not invalid token err : %v", err)
		return
	}

	// 클레임 추출
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorf("token.Claims fail")
		return
	}

	if val, exists := claims["sub"]; exists {
		bytes, _ := json.Marshal(val)
		facebookUser.UserID = string(bytes)
	}

	if val, exists := claims["email"]; exists {
		bytes, _ := json.Marshal(val)
		facebookUser.EA = string(bytes)
	}
}
