package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Application struct {
	AppID     int64  `json:"app_id" query:"app_id"`
	AppName   string `json:"app_name"`
	CompanyID int64  `json:"company_id"`
	Access    Access `json:"access"`
}

type Access struct {
	AccessID string `json:"access_id" validate:"required"`
	AccessPW string `json:"access_pw" validate:"required"`
}

type ResponseAppInfo struct {
	AppID     int64  `json:"app_id"`
	CompanyID int64  `json:"company_id"`
	AppName   string `json:"app_name"`
}

func NewApplication() *Application {
	return new(Application)
}

func (o *Access) CheckValidate() *base.BaseResponse {
	if len(o.AccessID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccessID)
	} else if len(o.AccessPW) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccessPW)
	}
	return nil
}
