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
	reqAuthAccountApp := new(context.ReqAuthAccountApplication)

	// Request json 파싱
	if err := c.Bind(reqAuthAccountApp); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// Request 유효성 체크
	if err := reqAuthAccountApp.CheckValidate(); err != nil {
		log.Error(err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostAccountLogin(c, reqAuthAccountApp)
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
