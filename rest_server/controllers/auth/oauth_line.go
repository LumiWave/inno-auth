package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// LINE 은 ID 토큰을 검증해 sub(사용자 ID)를 쓴다.
// access token verify 의 client_id 는 채널 ID 라서 모든 사용자가 같은 계정이 되므로 쓰면 안 됨.
var LineVerifyIDTokenEndpoint = "https://api.line.me/oauth2/v2.1/verify"

type LineIDToken struct {
	Sub string `json:"sub"`
	EA  string `json:"email"`
}

type OauthLine struct {
	SocialType int64
	ChannelID  string // 클라이언트(inno-game)가 ID 토큰을 발급받은 LINE 채널 ID. LINE 서버가 aud 와 대조함
}

func NewOauthLine() *OauthLine {
	return new(OauthLine)
}

func (o *OauthLine) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthLine) VerifySocialKey(idToken string) (string, string, error) {
	if o.ChannelID == "" {
		return "", "", errors.New("api_auth.line_channel_id 미설정")
	}

	resp, err := http.PostForm(LineVerifyIDTokenEndpoint, url.Values{
		"id_token":  {idToken},
		"client_id": {o.ChannelID},
	})
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// 만료·위조·다른 채널 토큰은 LINE 이 400 으로 거부
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("line id token verify: http %d", resp.StatusCode)
	}

	user := new(LineIDToken)
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return "", "", err
	}

	return user.Sub, user.EA, nil
}
