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

// Web 계정 로그인/가입
func (o *ExternalAPI) PostWebAccountLogin(c echo.Context) error {
	accountWeb := new(context.AccountWeb)

	// Request json 파싱
	if err := c.Bind(accountWeb); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// Request 유효성 체크
	if err := accountWeb.CheckValidate(); err != nil {
		log.Errorf("%v", err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWebAccountLogin(c, accountWeb)
}

// Web 계정 로그아웃
func (o *ExternalAPI) DelWebAccountLogout(c echo.Context) error {
	return commonapi.DelWebAccountLogout(c)
}

// Web 계정 로그인 정보 확인
func (o *ExternalAPI) PostWebAccountInfo(c echo.Context) error {
	reqAccountInfo := new(context.ReqAccountInfo)

	// Request json 파싱
	if err := c.Bind(reqAccountInfo); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	// Request 유효성 체크
	if err := reqAccountInfo.CheckValidate(); err != nil {
		log.Errorf("%v", err)
		return c.JSON(http.StatusOK, err)
	}
	return commonapi.PostWebAccountInfo(c, reqAccountInfo)
}
