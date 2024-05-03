// 회원 web account 로그인
package externalapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// Web 계정 로그인/가입
func (o *ExternalAPI) PostWebAccountLogin(c echo.Context) error {
	params := new(context.AccountWeb)

	// Request json 파싱
	if err := c.Bind(params); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}

	// Request 유효성 체크
	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWebAccountLogin(c, params, true)
}

// Web 계정 로그아웃
func (o *ExternalAPI) DelWebAccountLogout(c echo.Context) error {
	return commonapi.DelWebAccountLogout(c)
}

// Web 계정 로그인 정보 확인
func (o *ExternalAPI) PostWebAccountInfo(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	params := new(context.ReqAccountInfo)

	// Request 유효성 체크
	if err := params.CheckValidate(ctx); err != nil {
		log.Errorf("%v", err)
		return c.JSON(http.StatusOK, err)
	}

	return commonapi.PostWebAccountInfo(c, params)
}
