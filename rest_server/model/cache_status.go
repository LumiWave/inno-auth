package model

import (
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/gommon/log"
)

func MakeKeyStatus() string {
	return "INNO-STATUS-V2:STATUS"
}

func (o *DB) SetCacheStatus(params *context.StatusMain) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(MakeKeyStatus(), params, -1)
}

func (o *DB) GetCacheStatus() (*context.StatusMain, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	res := new(context.StatusMain)
	err := o.Cache.Get(MakeKeyStatus(), res)
	return res, err
}
