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
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	cKey := makeCacheKeyByAuth(payload.LoginType, tokenInfo.AccessUuid)
	err := o.Cache.Set(cKey, tokenInfo, time.Duration(accessExpiryPeriod))
	if err != nil {
		return err
	}

	cKey = makeCacheKeyByAuth(payload.LoginType, tokenInfo.RefreshUuid)
	err = o.Cache.Set(cKey, tokenInfo, time.Duration(refreshExpiryPeriod))
	if err != nil {
		return err
	}

	return nil
}

func (o *DB) GetJwtInfo(loginType context.LoginType, uuid string) (*context.JwtInfo, error) {
	cKey := makeCacheKeyByAuth(loginType, uuid)
	jwtInfo := new(context.JwtInfo)
	err := o.Cache.Get(cKey, jwtInfo)
	return jwtInfo, err
}

func (o *DB) DeleteJwtInfo(loginType context.LoginType, uuid string) error {
	cKey := makeCacheKeyByAuth(loginType, uuid)
	err := o.Cache.Del(cKey)
	return err
}

func makeCacheKeyByAuth(loginType context.LoginType, uuid string) string {
	return config.GetInstance().DBPrefix + ":INNO-AUTH-" + context.LoginTypeText[loginType] + ":" + uuid
}
