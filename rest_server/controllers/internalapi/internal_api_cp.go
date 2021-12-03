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
	cpInfo := context.NewCpInfo()
	if err := c.Bind(cpInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostCPRegister(c, cpInfo)
}

// CP사 탈퇴
func (o *InternalAPI) DelCPUnRegister(c echo.Context) error {
	cpInfo := context.NewCpInfo()
	if err := c.Bind(cpInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.DelCPUnRegister(c, cpInfo)
}
