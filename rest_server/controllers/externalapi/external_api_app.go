package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// App 로그인
func (o *ExternalAPI) PostAppLogin(c echo.Context) error {
	access := new(context.Access)

	if err := c.Bind(access); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := access.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostAppLogin(c, access)
}

// App 로그아웃
func (o *ExternalAPI) DelAppLogout(c echo.Context) error {
	return commonapi.DelAppLogout(c)
}
