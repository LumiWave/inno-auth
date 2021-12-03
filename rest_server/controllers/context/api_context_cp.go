package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Company struct {
	CompanyID   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
}

func NewCompany() *Company {
	return new(Company)
}

func (o *Company) CheckValidate() *base.BaseResponse {
	if len(o.CompanyName) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyCpName)
	}
	return nil
}
