// 회원 web account 로그인
package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) PostWebAccountLogin(c echo.Context) error {
	accountWeb := new(context.AccountWeb)

	// Request json 파싱
	if err := c.Bind(accountWeb); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// Request 유효성 체크
	if err := accountWeb.CheckValidate(); err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWebAccountLogin(c, accountWeb)
}