package model

import (
	"fmt"
	"time"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

func (o *DB) SetJwtInfoByCustomerUUID(tokenInfo *context.JwtInfo, payload *context.CustomerPayload) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	// Redis에 AccessToken 정보 등록
	cKey := MakeCacheKeyByUUID(payload.LoginType, context.AccessT, tokenInfo.AccessUuid)
	err := o.Cache.Set(cKey, tokenInfo, time.Duration(accessExpiryPeriod))
	if err != nil {
		return err
	}

	// Redis에 RefreshToken 정보 등록
	cKey = MakeCacheKeyByUUID(payload.LoginType, context.RefreshT, tokenInfo.RefreshUuid)
	err = o.Cache.Set(cKey, tokenInfo, time.Duration(refreshExpiryPeriod))
	if err != nil {
		return err
	}

	return nil
}

func (o *DB) GetJwtInfoByCustomerUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) (*context.JwtInfo, error) {
	cKey := MakeCacheKeyByUUID(loginType, tokenType, uuid)
	jwtInfo := new(context.JwtInfo)
	err := o.Cache.Get(cKey, jwtInfo)
	return jwtInfo, err
}

func (o *DB) DeleteJwtInfoByCustomerUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) error {
	cKey := MakeCacheKeyByUUID(loginType, tokenType, uuid)
	err := o.Cache.Del(cKey)
	return err
}

func MakeCacheKeyByCustomerUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) string {
	return fmt.Sprintf("%v:%v-%v:%v-%v", config.GetInstance().DBPrefix, "INNO-AUTH", context.LoginTypeText[loginType], context.TokenTypeText[tokenType], uuid)
}
