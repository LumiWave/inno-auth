package model

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

// 로그인 성공시 정보 추가
func (o *DB) SetJwtInfo(tokenInfo *context.JwtInfo, appInfo *context.AppInfo) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}
	conf := config.GetInstance()

	cKey := makeCacheKeyByAuth(tokenInfo.AccessUuid)
	err := o.Cache.Set(cKey, appInfo, time.Duration(conf.Auth.TokenExpiryPeriod*int64(time.Minute)))
	if err != nil {
		return err
	}

	cKey = makeCacheKeyByAuth(tokenInfo.RefreshUuid)
	err = o.Cache.Set(cKey, appInfo, time.Duration(conf.Auth.TokenExpiryPeriod*int64(time.Minute)))
	if err != nil {
		return err
	}

	return nil
}

func (o *DB) GetJwtInfo(uuid string) (*context.AppInfo, error) {
	cKey := makeCacheKeyByAuth(uuid)
	appInfo := new(context.AppInfo)
	err := o.Cache.Get(cKey, appInfo)
	return appInfo, err
}

func (o *DB) DeleteJwtInfo(uuid string) error {
	cKey := makeCacheKeyByAuth(uuid)
	err := o.Cache.Del(cKey)
	return err
}

func makeCacheKeyByAuth(id string) string {
	return config.GetInstance().DBPrefix + ":INNO-AUTH-APP:" + id
}
