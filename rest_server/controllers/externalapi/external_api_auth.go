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
