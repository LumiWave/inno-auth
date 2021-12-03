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
	app := context.NewApplication()
	if err := c.Bind(app); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostAppRegister(c, app)
}

// App 삭제
func (o *InternalAPI) DelAppUnRegister(c echo.Context) error {
	app := context.NewApplication()
	if err := c.Bind(app); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.DelAppUnRegister(c, app)
}
