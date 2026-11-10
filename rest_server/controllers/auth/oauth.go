package auth

import "github.com/LumiWave/inno-auth/rest_server/model"

const (
	SocialType_Google   = 1
	SocialType_Facebook = 2
	SocialType_Inno     = 3
	SocialType_LINE     = 4
	SocialType_Telegram = 5
	SocialType_Apple    = 6
)

type SocialAuth interface {
	GetSocialType() int64
	VerifySocialKey(string) (string, string, error)
}

// Telegram 은 차단: oauth_telegram.go 가 봇 토큰 없이(빈 문자열) 서명을 검증해 누구나 로그인 값을 위조할 수 있음.
// 봇 토큰 검증 + auth_date 만료 검사를 고치기 전에는 다시 넣지 말 것.
func CheckValidateExternal(socialType int64) bool {
	if socialType == SocialType_Google ||
		socialType == SocialType_Facebook ||
		socialType == SocialType_LINE ||
		socialType == SocialType_Apple {
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
		case social.SocialType == SocialType_LINE:
			isocial := OauthLine{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		case social.SocialType == SocialType_Telegram:
			isocial := OauthTelegram{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		case social.SocialType == SocialType_Apple:
			isocial := OauthApple{
				SocialType: social.SocialType,
			}
			socialAuths[social.SocialType] = &isocial
		}
	}

	iAuth.SocialAuths = socialAuths
}
