package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAppRegister(c echo.Context, appInfo *context.AppInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// App 이름 빈문자열 체크
	if err := appInfo.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// App 이름 중복 체크
	if respAppInfo, err := model.GetDB().SelectGetAppInfoByAppName(appInfo.AppName); err == nil {
		if len(respAppInfo.AppName) > 0 {
			log.Error("PostAppRegister exists app_name=", appInfo.AppName, " errorCode:", resultcode.Result_Auth_ExistsAppName)
			resp.SetReturn(resultcode.Result_Auth_ExistsAppName)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// cp_idx 존재 여부 확인
	if value, err := model.GetDB().SelectGetCpInfoByIdx(appInfo.CpIdx); err == nil {
		if len(value.CpName) == 0 {
			log.Error("PostCPRegister Not Exists cp_idx=", appInfo.CpIdx, " errorCode:", resultcode.Result_Auth_NotFoundCpIdx)
			resp.SetReturn(resultcode.Result_Auth_NotFoundCpName)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 테이블에 신규 row 생성
	if err := model.GetDB().InsertApp(appInfo); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return c.JSON(http.StatusOK, resp)
}

func DelAppUnRegister(c echo.Context, appInfo *context.AppInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// App 이름 빈문자열 체크
	if err := appInfo.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 테이블 row 삭제
	if err := model.GetDB().DeleteApp(appInfo); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}
	return c.JSON(http.StatusOK, resp)
}

func GetAppExists(c echo.Context, appInfo *context.AppInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if value, err := model.GetDB().SelectGetAppInfoByAppName(appInfo.AppName); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if len(value.AppName) != 0 {
			resp.Value = value
		} else {
			resp.SetReturn(resultcode.Result_Auth_NotFoundAppName)
		}
	}
	return c.JSON(http.StatusOK, resp)
}

func PostAppLogin(c echo.Context, reqAppInfo *context.RequestAppLoginInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 가입 정보 확인

	// 2. redis duplicate check
	// redis에 기존 정보가 있다면 기존에 발급된 토큰으로 응답한다.
	appInfo := new(context.AppInfo)
	appInfo.Account = reqAppInfo.Account

	// 3. create Auth Token
	if jwtInfoValue, err := auth.GetIAuth().MakeToken(appInfo); err != nil {
		resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
	} else {
		resp.Value = jwtInfoValue
	}

	return c.JSON(http.StatusOK, resp)
}

func DelAppLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	auth.GetIAuth().DeleteJwtInfo(ctx.Uuid)

	return c.JSON(http.StatusOK, resp)
}
