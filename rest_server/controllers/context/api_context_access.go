package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Access struct {
	AccessID string `json:"access_id" validate:"required"`
	AccessPW string `json:"access_pw" validate:"required"`
}

func (o *Access) CheckValidate() *base.BaseResponse {
	if len(o.AccessID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccessID)
	} else if len(o.AccessPW) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccessPW)
	}
	return nil
}
