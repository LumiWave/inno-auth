package auth

type OauthFacebook struct {
	SocialType int64
}

func NewOauthFacebook() *OauthFacebook {
	return new(OauthFacebook)
}

func (o *OauthFacebook) GetSocialType() int64 {
	return o.SocialType
}

func (o *OauthFacebook) VerifySocialKey(socialKey string) error {

	return nil
}

func (o *OauthFacebook) MakeInnoId() (string, error) {
	return "", nil
}
