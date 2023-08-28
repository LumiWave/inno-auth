package context

import "github.com/ONBUFF-IP-TOKEN/baseapp/base"

type ReqIPCheck struct {
	Ip string `json:"ip"`
}

type RespIPCheck struct {
	Country     string `json:"country"`
	AllowAccess bool   `json:"allow_access"`
	SwapEnable  bool   `json:"swap_enable"`
}

// swap enable 설정
type ReqSwapEnable struct {
	SwapEnable bool `json:"swap_enable"`
}

func (o *ReqSwapEnable) CheckValidate() *base.BaseResponse {
	return nil
}
