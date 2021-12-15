package auth

import (
	"encoding/json"
	"io/ioutil"

	contextR "context"

	"golang.org/x/oauth2"
)

const (
	UserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type GoogleUser struct {
	Email  string `json:"email"`
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
	userInfoResp, err := client.Get(UserInfoAPIEndpoint)
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

	return googleUser.UserID, googleUser.Email, nil
}

func (o *OauthGoogle) MakeInnoId() (string, error) {
	return "", nil
}
