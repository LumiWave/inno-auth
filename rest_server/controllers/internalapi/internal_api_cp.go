package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func ParseCPInfo(c echo.Context) (*context.CpInfo, error) {
	params := context.NewCpInfo()
	if err := c.Bind(params); err != nil {
		log.Error(err)
		return params, base.BaseJSONInternalServerError(c, err)
	}
	return params, nil
}

// CP사 신규 가입
func (o *InternalAPI) PostCPRegister(c echo.Context) error {
	params, err := ParseCPInfo(c)
	if err != nil {
		return err
	}
	context.MakeDt(&params.CreateDt)

	return commonapi.PostCPRegister(c, params)
}

// CP사 탈퇴
func (o *InternalAPI) DelCPUnRegister(c echo.Context) error {
	params, err := ParseCPInfo(c)
	if err != nil {
		return err
	}
	return commonapi.DelCPUnRegister(c, params)
}
