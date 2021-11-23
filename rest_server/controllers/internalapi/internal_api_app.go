package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// App 신규 추가
func (o *InternalAPI) PostAppRegister(c echo.Context) error {
	params := context.NewAppInfo()

	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	context.MakeDt(&params.CreateDt)

	return commonapi.PostAppRegister(c, params)
}
