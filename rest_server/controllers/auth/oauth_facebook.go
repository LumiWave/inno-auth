package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	FacebookUserInfoAPIEndpoint = "https://graph.facebook.com/me?fields=email&access_token="
)

type FacebookUser struct {
	Email  string `json:"email"`
	UserID string `json:"id"`
}
type OauthFacebook struct {
	SocialType int
}

func NewOauthFacebook() *OauthFacebook {
	return new(OauthFacebook)
}

func (o *OauthFacebook) GetSocialType() int {
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

	return facebookUser.UserID, facebookUser.Email, nil
}
