package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/model"
)

// 서명이 맞는 토큰이라도, 서버가 실제로 발급해서 Redis 에 보관 중인 토큰과 같아야 인정한다.
// 서명키를 아는 사람은 서명이 맞는 토큰을 얼마든지 만들 수 있으므로 서명만으로는 신원을 믿을 수 없다.
// 예전에는 "그 계정의 세션이 있는가"만 봤기 때문에, 남의 InnoUID 나 au_id 를 적은 위조 토큰이 통과했다.

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionMismatch = errors.New("token is not the issued one")

// TokenHash 는 이전 토큰 목록에 보관하는 토큰의 해시다.
func TokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func sameToken(a, b string) bool {
	return len(a) != 0 && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func issuedToken(jwtInfo *context.JwtInfo, tokenType context.TokenType) string {
	if jwtInfo == nil {
		return ""
	}
	if tokenType == context.RefreshT {
		return jwtInfo.RefreshToken
	}
	return jwtInfo.AccessToken
}

// matchSession 은 제시된 토큰이 보관 중인 세션의 토큰이거나, 갱신 전에 발급했던 토큰인지 확인한다.
func matchSession(jwtInfo *context.JwtInfo, tokenType context.TokenType, token string, prevHashes []string) bool {
	if sameToken(issuedToken(jwtInfo, tokenType), token) {
		return true
	}
	hash := TokenHash(token)
	for _, prev := range prevHashes {
		if sameToken(prev, hash) {
			return true
		}
	}
	return false
}

// CheckSession 은 제시된 토큰이 서버가 발급해 보관 중인 토큰인지 확인한다.
// uuid 는 토큰의 access_uuid/refresh_uuid, innoUID 는 웹 회원 토큰의 inno_uid 다.
func (o *IAuth) CheckSession(loginType context.LoginType, tokenType context.TokenType, uuid, innoUID, token string) error {
	switch loginType {
	case context.AppLogin:
		jwtInfo, err := o.GetJwtInfoByUUID(loginType, tokenType, uuid)
		if err != nil {
			return ErrSessionNotFound
		}
		if !matchSession(jwtInfo, tokenType, token, nil) {
			return ErrSessionMismatch
		}
	case context.CustomerLogin:
		jwtInfo, err := o.GetJwtInfoByCustomerUUID(loginType, tokenType, uuid)
		if err != nil {
			return ErrSessionNotFound
		}
		if !matchSession(jwtInfo, tokenType, token, nil) {
			return ErrSessionMismatch
		}
	case context.WebAccountLogin:
		jwtInfo, err := o.GetJwtInfoByInnoUID(loginType, tokenType, innoUID)
		if err != nil {
			return ErrSessionNotFound
		}
		// 다른 기기에서 갱신하기 전에 받은 토큰도 로그아웃 전까지는 인정한다.
		if !matchSession(jwtInfo, tokenType, token, model.GetDB().GetPrevTokenHashes(loginType, innoUID)) {
			return ErrSessionMismatch
		}
	default:
		// 발급한 적 없는 로그인 종류의 토큰
		return ErrSessionNotFound
	}
	return nil
}
