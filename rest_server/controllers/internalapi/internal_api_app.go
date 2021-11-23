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
	appInfo := context.NewAppInfo()

	if err := c.Bind(appInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	context.MakeDt(&appInfo.CreateDt)

	return commonapi.PostAppRegister(c, appInfo)
}
