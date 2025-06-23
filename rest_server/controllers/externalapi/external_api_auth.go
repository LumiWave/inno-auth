package externalapi

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) GetTokenVerify(c echo.Context) error {
	return commonapi.GetTokenVerify(c)
}

func (o *ExternalAPI) PostTokenRenew(c echo.Context) error {
	renewTokenRequest := new(context.RenewTokenRequest)

	if err := c.Bind(renewTokenRequest); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostTokenRenew(c, renewTokenRequest)
}

func (o *ExternalAPI) PostIPAccessAllow(c echo.Context) error {
	reqIpCheck := new(context.ReqIPCheck)
	// Request json 파싱
	if err := c.Bind(reqIpCheck); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostIPAccessAllow(c, reqIpCheck)
}

func (o *ExternalAPI) GetPermissionAvailable(c echo.Context) error {

	reqPA := new(context.ReqPermissionAvailable)
	if err := c.Bind(reqPA); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.GetPermissionAvailable(c, reqPA)
}
