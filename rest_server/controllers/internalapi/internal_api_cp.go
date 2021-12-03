package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// CP사 신규 가입
func (o *InternalAPI) PostCPRegister(c echo.Context) error {
	company := context.NewCompany()
	if err := c.Bind(company); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostCPRegister(c, company)
}

// CP사 탈퇴
func (o *InternalAPI) DelCPUnRegister(c echo.Context) error {
	company := context.NewCompany()
	if err := c.Bind(company); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.DelCPUnRegister(c, company)
}
