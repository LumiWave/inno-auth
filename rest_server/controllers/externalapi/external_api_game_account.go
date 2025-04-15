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
func (o *ExternalAPI) PostGameAccountLogin(c echo.Context) error {
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

	return commonapi.PostGameAccountLogin(c, params, true)
}

// Web Account 인증(app에서 직접 가입/로그인)
// func (o *ExternalAPI) PostGameAccountLoginOnce(c echo.Context) error {
// 	params := new(context.AccountWeb)

// 	// Request json 파싱
// 	if err := c.Bind(params); err != nil {
// 		log.Errorf("%v", err)
// 		return base.BaseJSONInternalServerError(c, err)
// 	}

// 	// Request 유효성 체크
// 	if err := params.CheckValidate(); err != nil {
// 		return c.JSON(http.StatusOK, err)
// 	}

// 	return commonapi.PostWebAccountLoginOnce(c, params, true)
// }

// Web 계정 로그아웃
func (o *ExternalAPI) DelGameAccountLogout(c echo.Context) error {
	return commonapi.DelGameAccountLogout(c)
}

// Web 계정 로그인 정보 확인
// func (o *ExternalAPI) PostGameAccountInfo(c echo.Context) error {
// 	ctx := base.GetContext(c).(*context.InnoAuthContext)
// 	params := new(context.ReqAccountInfo)

// 	// Request 유효성 체크
// 	if err := params.CheckValidate(ctx); err != nil {
// 		log.Errorf("%v", err)
// 		return c.JSON(http.StatusOK, err)
// 	}

// 	return commonapi.PostWebAccountInfo(c, params)
// }
