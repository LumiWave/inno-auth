package model

import (
	"fmt"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

func (o *DB) SaveRedisInfoByInnoUIDGame(jwtInfo *context.JwtInfo, payload *context.Payload, tokenType context.TokenType, expiryPeriod int64) error {
	cKey := MakeCacheKeyByInnoUIDGame(payload.LoginType, tokenType, payload.InnoUID)
	err := o.Cache.Set(cKey, jwtInfo, time.Duration(expiryPeriod))
	if err != nil {
		return err
	}
	return nil
}

func (o *DB) SetJwtInfoByInnoUIDGame(jwtInfo *context.JwtInfo, payload *context.Payload) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	// Redis에 AccessToken 정보 등록
	if err := o.SaveRedisInfoByInnoUIDGame(jwtInfo, payload, context.AccessT, accessExpiryPeriod); err != nil {
		log.Errorf("%v", err)
		return err
	}

	// Redis에 RefreshToken 정보 등록
	if err := o.SaveRedisInfoByInnoUIDGame(jwtInfo, payload, context.RefreshT, refreshExpiryPeriod); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func (o *DB) GetJwtInfoByInnoUIDGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) (*context.JwtInfo, error) {
	cKey := MakeCacheKeyByInnoUIDGame(loginType, tokenType, innoUID)
	jwtInfo := new(context.JwtInfo)
	err := o.Cache.Get(cKey, jwtInfo)
	return jwtInfo, err
}

func (o *DB) DeleteJwtInfoByInnoUIDGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	cKey := MakeCacheKeyByInnoUIDGame(loginType, tokenType, innoUID)
	err := o.Cache.Del(cKey)
	return err
}

func MakeCacheKeyByInnoUIDGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) string {
	return fmt.Sprintf("%v:%v-%v:%v-%v", config.GetInstance().DBPrefix, "INNO-AUTH-GAME", context.LoginTypeText[loginType], context.TokenTypeText[tokenType], innoUID)
}
