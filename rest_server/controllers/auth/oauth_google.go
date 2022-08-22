package auth

import (
	"encoding/json"
	"io/ioutil"

	contextR "context"

	"golang.org/x/oauth2"
)

const (
	GoogleUserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type GoogleUser struct {
	EA     string `json:"email"`
	UserID string `json:"sub"`
}

type OauthGoogle struct {
	SocialType int64
}

func NewOauthGoogle() *OauthGoogle {
	return new(OauthGoogle)
}

func (o *OauthGoogle) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthGoogle) VerifySocialKey(socialKey string) (string, string, error) {
	token := &oauth2.Token{
		AccessToken: socialKey,
		TokenType:   "Bearer",
	}
	var OAuthConf *oauth2.Config
	client := OAuthConf.Client(contextR.Background(), token)
	userInfoResp, err := client.Get(GoogleUserInfoAPIEndpoint)
	if err != nil || userInfoResp.StatusCode != 200 {
		return "", "", err
	}
	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		return "", "", err
	}
	googleUser := new(GoogleUser)
	json.Unmarshal(userInfo, &googleUser)

	return googleUser.UserID, googleUser.EA, nil
}
