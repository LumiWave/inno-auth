package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type CpInfo struct {
	CompanyID   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
}

func NewCpInfo() *CpInfo {
	return new(CpInfo)
}

func (o *CpInfo) CheckValidate() *base.BaseResponse {
	if len(o.CompanyName) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyCpName)
	}
	return nil
}
