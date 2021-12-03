package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Application struct {
	AppID     int        `json:"app_id" query:"app_id"`
	AppName   string     `json:"app_name"`
	CompanyID int        `json:"company_id"`
	Access    AccessInfo `json:"access_info"`
}

type RequestAppLoginInfo struct {
	Access AccessInfo `json:"access_info" validate:"required"`
}

type ResponseAppInfo struct {
	AppID     int    `json:"app_id"`
	CompanyID int    `json:"company_id"`
	AppName   string `json:"app_name"`
}

func NewApplication() *Application {
	return new(Application)
}

func (o *Application) CheckValidate() *base.BaseResponse {
	if len(o.AppName) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAppName)
	}
	return nil
}
