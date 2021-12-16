package auth

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

	return "", "", nil
}
