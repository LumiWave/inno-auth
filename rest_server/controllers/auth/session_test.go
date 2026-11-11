package auth

import (
	"testing"

	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

func TestMatchSession(t *testing.T) {
	issued := &context.JwtInfo{AccessToken: "header.payload-real.sig", RefreshToken: "header.refresh-real.sig"}
	oldAccess := "header.payload-old.sig"
	prev := []string{TokenHash(oldAccess), TokenHash("header.refresh-old.sig")}

	cases := []struct {
		name      string
		jwtInfo   *context.JwtInfo
		tokenType context.TokenType
		token     string
		prev      []string
		want      bool
	}{
		{"발급한 access 토큰", issued, context.AccessT, "header.payload-real.sig", nil, true},
		{"발급한 refresh 토큰", issued, context.RefreshT, "header.refresh-real.sig", nil, true},
		{"같은 세션 번호를 넣어 위조한 토큰 (내용이 다름)", issued, context.AccessT, "header.payload-forged.sig", nil, false},
		{"access 자리에 refresh 토큰", issued, context.AccessT, "header.refresh-real.sig", nil, false},
		{"갱신 전에 발급했던 토큰", issued, context.AccessT, oldAccess, prev, true},
		{"갱신 전 목록에도 없는 토큰", issued, context.AccessT, "header.payload-forged.sig", prev, false},
		{"로그아웃해서 이전 목록이 지워진 뒤의 옛 토큰", issued, context.AccessT, oldAccess, nil, false},
		{"빈 토큰", issued, context.AccessT, "", prev, false},
		{"세션에 토큰이 비어 있는 경우 빈 토큰", &context.JwtInfo{}, context.AccessT, "", nil, false},
		{"세션 없음", nil, context.AccessT, "header.payload-real.sig", nil, false},
	}
	for _, c := range cases {
		if got := matchSession(c.jwtInfo, c.tokenType, c.token, c.prev); got != c.want {
			t.Errorf("%s : %v, 기대 %v", c.name, got, c.want)
		}
	}
}

func TestTokenHash(t *testing.T) {
	a, b := TokenHash("token-a"), TokenHash("token-b")
	if a == b {
		t.Error("서로 다른 토큰의 해시가 같음")
	}
	if a != TokenHash("token-a") {
		t.Error("같은 토큰의 해시가 다름")
	}
	if len(a) != 64 {
		t.Errorf("해시 길이 %d, 기대 64", len(a))
	}
}
