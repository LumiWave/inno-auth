package auth

type OauthGoogle struct {
	SocialType int64
}

func NewOauthGoogle() *OauthGoogle {
	return new(OauthGoogle)
}

func (o OauthGoogle) GetSocialType() int64 {
	return o.SocialType
}

func (o OauthGoogle) VerifySocialKey(socialKey string) error {

	return nil
}

func (o OauthGoogle) MakeInnoId() (string, error) {
	return "", nil
}
