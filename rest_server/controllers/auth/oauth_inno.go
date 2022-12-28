package auth

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

type AIUsaa struct {
	EA     string `json:"email"`
	UserID string `json:"sub"`
}

type OauthAI struct {
	SocialType int64
}

func NewOauthAI() *OauthAI {
	return new(OauthAI)
}

func (o *OauthAI) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthAI) VerifySocialKey(socialKey string) (string, string, error) {
	conf := config.GetInstance()
	userID := inno.AESDecrypt(socialKey, []byte(conf.Secret.Key), []byte(conf.Secret.Iv))

	return userID, "pick@onbuff.com", nil
}
