package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type AppInfo struct {
	Idx      int64       `json:"idx" query:"idx"`
	AppName  string      `json:"app_name" query:"app_name" validate:"required"`
	CpIdx    int64       `json:"cp_idx" query:"cp_idx" validate:"required"`
	Account  AccountInfo `json:"account_info" query:"app_name" validate:"required"`
	Token    JwtInfo
	CreateDt int64 `json:"create_dt"`
}

type ResponseAppInfo struct {
	Idx     int64  `json:"idx"`
	CpIdx   int64  `json:"cp_idx"`
	AppName string `json:"app_name"`
}

func NewAppInfo() *AppInfo {
	return new(AppInfo)
}

func NewRespAppInfo() *ResponseAppInfo {
	return new(ResponseAppInfo)
}

func (o *AppInfo) CheckValidate() *base.BaseResponse {
	if len(o.AppName) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAppName)
	}
	return nil
}
