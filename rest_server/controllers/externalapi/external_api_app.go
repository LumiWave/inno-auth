package externalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// App 로그인
func (o *ExternalAPI) PostAppLogin(c echo.Context) error {
	access := new(context.Access)

	if err := c.Bind(access); err != nil {
		log.Errorf("%v", err)
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
