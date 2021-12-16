package auth

import "github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"

const (
	SocialType_Google   = 1
	SocialType_Facebook = 2
)

type SocialAuth interface {
	GetSocialType() int
	VerifySocialKey(string) (string, string, error)
}

func MakeSocialAuths(iAuth *IAuth) {
	socialAuths := make(map[int]SocialAuth)

	for _, social := range model.GetDB().Socials {
		switch {
		case social.SocialType == SocialType_Google:
			isocial := OauthGoogle{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		case social.SocialType == SocialType_Facebook:
			isocial := OauthFacebook{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		}
	}

	iAuth.SocialAuths = socialAuths
}
