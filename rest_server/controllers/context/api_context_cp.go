package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type CpInfo struct {
	Idx      int64  `json:"idx" query:"idx"`
	CpName   string `json:"cp_name" query:"cp_name" validate:"required"`
	Token    TokenInfo
	CreateDt int64 `json:"create_dt"`
}

func NewCpInfo() *CpInfo {
	return new(CpInfo)
}

func (o *CpInfo) CheckValidate() *base.BaseResponse {
	if len(o.CpName) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyCpName)
	}
	return nil
}
