package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	LineVierfyAccessTokenValidityEndpoint = "https://api.line.me/oauth2/v2.1/verify?access_token="
)

type LineUser struct {
	Scope     string `json:"scope"`
	ClientID  string `json:"client_id"`
	ExpiresIn int64  `json:"expires_in"`
	EA        string `json:"email"`
}

type OauthLine struct {
	SocialType int64
}

func NewOauthLine() *OauthLine {
	return new(OauthLine)
}

func (o *OauthLine) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthLine) VerifySocialKey(socialKey string) (string, string, error) {
	userInfoResp, err := http.Get(LineVierfyAccessTokenValidityEndpoint + url.QueryEscape(socialKey))
	if err != nil {
		return "", "", err
	}
	defer userInfoResp.Body.Close()

	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		return "", "", err
	}

	user := new(LineUser)
	json.Unmarshal(userInfo, &user)

	return user.ClientID, user.EA, nil
}
