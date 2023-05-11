package externalapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 고객사 로그인
func (o *ExternalAPI) PostCustomerLogin(c echo.Context) error {
	access := new(context.CustomerAccess)

	if err := c.Bind(access); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}

	if err := access.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostCustomerLogin(c, access)
}

// 고객사 로그아웃
func (o *ExternalAPI) DelCustomerLogout(c echo.Context) error {
	return commonapi.DelCustomerLogout(c)
}

// 고객사 액세스 토큰 검증
func (o *ExternalAPI) GetCustomerTokenVerify(c echo.Context) error {
	return commonapi.GetCustomerTokenVerify(c)
}

// 고객사 액세스 토큰 갱신/재발급
func (o *ExternalAPI) PostCustomerTokenRenew(c echo.Context) error {
	renewTokenRequest := new(context.RenewTokenRequest)

	if err := c.Bind(renewTokenRequest); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.PostCustomerTokenRenew(c, renewTokenRequest)
}
