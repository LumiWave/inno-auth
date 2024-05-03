package auth

import "github.com/LumiWave/inno-auth/rest_server/model"

const (
	SocialType_Google   = 1
	SocialType_Facebook = 2
	SocialType_Inno     = 3
)

type SocialAuth interface {
	GetSocialType() int64
	VerifySocialKey(string) (string, string, error)
}

func CheckValidateExternal(socialType int64) bool {
	if socialType == SocialType_Google || socialType == SocialType_Facebook {
		return true
	}
	return false
}

func CheckValidateInternal(socialType int64) bool {
	return socialType == SocialType_Inno
}

func MakeSocialAuths(iAuth *IAuth) {
	socialAuths := make(map[int64]SocialAuth)

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
		case social.SocialType == SocialType_Inno:
			isocial := OauthAI{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		}
	}

	iAuth.SocialAuths = socialAuths
}
