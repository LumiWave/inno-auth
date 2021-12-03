package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 회원 로그인
func (o *ExternalAPI) PostAccountLogin(c echo.Context) error {
	//ctx := base.GetContext(c).
	reqAccountAuth := new(context.RequestAccountAuth)

	if err := c.Bind(reqAccountAuth); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := reqAccountAuth.CheckValidate(); err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostAccountLogin(c, reqAccountAuth)
}

// App 존재 여부 확인
func (o *ExternalAPI) GetAccountExists(c echo.Context) error {
	account := context.NewAccount()

	if err := c.Bind(account); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.GetAccountExists(c, account)
}
