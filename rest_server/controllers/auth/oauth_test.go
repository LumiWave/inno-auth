package auth

import "testing"

func TestCheckValidateExternal(t *testing.T) {
	if CheckValidateExternal(SocialType_Telegram) {
		t.Fatal("Telegram 은 외부 로그인에서 차단돼야 함 (빈 봇 토큰 검증 취약점)")
	}
	if CheckValidateExternal(SocialType_Inno) {
		t.Fatal("Inno 는 내부 API 전용")
	}
	for _, s := range []int64{SocialType_Google, SocialType_Facebook, SocialType_LINE, SocialType_Apple} {
		if !CheckValidateExternal(s) {
			t.Fatalf("social type %d 는 외부 로그인 허용돼야 함", s)
		}
	}
}
