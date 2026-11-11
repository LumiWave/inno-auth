package model

import (
	"fmt"
	"time"

	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

// 토큰을 갱신하면 이전 토큰은 세션에서 사라진다. 같은 계정으로 여러 기기에서 쓰는 경우
// 다른 기기는 이전 토큰을 들고 있으므로, 갱신 전 토큰의 해시를 따로 보관해 두고 계속 인정한다.
// 해시만 보관하므로 Redis 가 노출돼도 토큰 자체는 복원되지 않는다.

// 계정 하나에 보관하는 이전 토큰 해시 수 (access, refresh 각각 포함)
const maxPrevTokenHashes = 20

type PrevTokens struct {
	Hashes []string `json:"hashes"`
}

func MakeCacheKeyPrevTokens(loginType context.LoginType, innoUID string) string {
	return fmt.Sprintf("%v:%v-%v:%v-%v", config.GetInstance().DBPrefix, "INNO-AUTH", context.LoginTypeText[loginType], "PREV", innoUID)
}

// 이전 토큰 해시 목록 조회. 없으면 빈 목록.
func (o *DB) GetPrevTokenHashes(loginType context.LoginType, innoUID string) []string {
	prev := new(PrevTokens)
	if err := o.Cache.Get(MakeCacheKeyPrevTokens(loginType, innoUID), prev); err != nil {
		return nil
	}
	return prev.Hashes
}

// 이전 토큰 해시 추가. refresh 토큰 유효 기간만큼 보관한다.
func (o *DB) AddPrevTokenHashes(loginType context.LoginType, innoUID string, hashes ...string) error {
	list := append(o.GetPrevTokenHashes(loginType, innoUID), hashes...)
	if len(list) > maxPrevTokenHashes {
		list = list[len(list)-maxPrevTokenHashes:]
	}
	_, refreshExpiryPeriod := context.GetTokenExpiryperiod(loginType)
	return o.Cache.Set(MakeCacheKeyPrevTokens(loginType, innoUID), &PrevTokens{Hashes: list}, time.Duration(refreshExpiryPeriod))
}

// 로그아웃하면 이전 토큰도 전부 무효로 만든다.
func (o *DB) DeletePrevTokens(loginType context.LoginType, innoUID string) error {
	return o.Cache.Del(MakeCacheKeyPrevTokens(loginType, innoUID))
}
