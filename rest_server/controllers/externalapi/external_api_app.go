package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// App 로그인
func (o *ExternalAPI) PostAppLogin(c echo.Context) error {
	appInfo := context.NewRequestAppLoginInfo()

	if err := c.Bind(appInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostAppLogin(c, appInfo)
}

// App 로그아웃
func (o *ExternalAPI) DelAppLogout(c echo.Context) error {
	return commonapi.DelAppLogout(c)
}

// App 존재 여부 확인
func (o *ExternalAPI) GetAppExists(c echo.Context) error {
	appInfo := context.NewAppInfo()

	if err := c.Bind(appInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.GetAppExists(c, appInfo)
}
