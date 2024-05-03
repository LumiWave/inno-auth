package externalapi

import (
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) PostSuiProver(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	params := new(context.ReqProve)

	if err := ctx.EchoContext.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostSuiProver(ctx, params)
}
