package model

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

// 로그인 성공시 정보 추가
func (o *DB) SetJwtInfo(tokenInfo *context.JwtInfo, payload *context.Payload) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	conf := config.GetInstance()

	cKey := makeCacheKeyByAuth(payload, tokenInfo.AccessUuid)
	err := o.Cache.Set(cKey, tokenInfo, time.Duration(conf.Auth.AccessTokenExpiryPeriod*int64(time.Minute)))
	if err != nil {
		return err
	}

	cKey = makeCacheKeyByAuth(payload, tokenInfo.RefreshUuid)
	err = o.Cache.Set(cKey, tokenInfo, time.Duration(conf.Auth.RefreshTokenExpiryPeriod*int64(time.Minute)))
	if err != nil {
		return err
	}

	return nil
}

func (o *DB) GetJwtInfo(payload *context.Payload) (*context.JwtInfo, error) {
	cKey := makeCacheKeyByAuth(payload, payload.Uuid)
	jwtInfo := new(context.JwtInfo)
	err := o.Cache.Get(cKey, jwtInfo)
	return jwtInfo, err
}

func (o *DB) DeleteJwtInfo(payload *context.Payload) error {
	cKey := makeCacheKeyByAuth(payload, payload.Uuid)
	err := o.Cache.Del(cKey)
	return err
}

func makeCacheKeyByAuth(payload *context.Payload, uuid string) string {
	return config.GetInstance().DBPrefix + ":INNO-AUTH-" + context.LoginTypeText[payload.LoginType] + ":" + uuid
}
