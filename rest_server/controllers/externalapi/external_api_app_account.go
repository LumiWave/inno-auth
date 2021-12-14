package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 회원 app account 로그인
func (o *ExternalAPI) PostAppAccountLogin(c echo.Context) error {
	account := context.NewAccount()

	// Request json 파싱
	if err := c.Bind(account); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// Request 유효성 체크
	if err := account.CheckValidate(); err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostAppAccountLogin(c, account)
}
