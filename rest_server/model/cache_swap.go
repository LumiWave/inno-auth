package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

func (o *DB) GetSwapEnable() (bool, error) {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	bEnable := false
	err := o.Cache.Get(MakeCacheKeySwapEnable(), &bEnable)
	return bEnable, err
}

func (o *DB) SetSwapEnable(bEnable bool) error {
	if !o.Cache.Enable() {
		log.Warnf("redis disable")
	}

	return o.Cache.Set(MakeCacheKeySwapEnable(), bEnable, -1)
}

func MakeCacheKeySwapEnable() string {
	return fmt.Sprintf("%v:%v", config.GetInstance().DBPrefix, "SWAP-ENABLE")
}
