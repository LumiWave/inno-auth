package model

import (
	"fmt"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

func (o *DB) SaveRedisInfoByInnoUID(jwtInfo *context.JwtInfo, payload *context.Payload, tokenType context.TokenType, expiryPeriod int64) error {
	cKey := MakeCacheKeyByInnoUID(payload.LoginType, tokenType, payload.InnoUID)
	err := o.Cache.Set(cKey, jwtInfo, time.Duration(expiryPeriod))
	if err != nil {
		return err
	}
	return nil
}

func (o *DB) SetJwtInfoByInnoUID(jwtInfo *context.JwtInfo, payload *context.Payload) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	// Redis에 AccessToken 정보 등록
	if err := o.SaveRedisInfoByInnoUID(jwtInfo, payload, context.AccessT, accessExpiryPeriod); err != nil {
		log.Errorf("%v", err)
		return err
	}

	// Redis에 RefreshToken 정보 등록
	if err := o.SaveRedisInfoByInnoUID(jwtInfo, payload, context.RefreshT, refreshExpiryPeriod); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

func (o *DB) GetJwtInfoByInnoUID(loginType context.LoginType, tokenType context.TokenType, innoUID string) (*context.JwtInfo, error) {
	cKey := MakeCacheKeyByInnoUID(loginType, tokenType, innoUID)
	jwtInfo := new(context.JwtInfo)
	err := o.Cache.Get(cKey, jwtInfo)
	return jwtInfo, err
}

func (o *DB) DeleteJwtInfoByInnoUID(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	cKey := MakeCacheKeyByInnoUID(loginType, tokenType, innoUID)
	err := o.Cache.Del(cKey)
	return err
}

func MakeCacheKeyByInnoUID(loginType context.LoginType, tokenType context.TokenType, innoUID string) string {
	return fmt.Sprintf("%v:%v-%v:%v-%v", config.GetInstance().DBPrefix, "INNO-AUTH", context.LoginTypeText[loginType], context.TokenTypeText[tokenType], innoUID)
}
