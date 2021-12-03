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

	// companyID 존재 여부 확인
	if value, err := model.GetDB().SelectGetCpInfoByIdx(appInfo.CompanyID); err == nil {
		if len(value.CompanyName) == 0 {
			log.Error("PostCPRegister Not Exists CompanyID=", appInfo.CompanyID, " errorCode:", resultcode.Result_Auth_NotFoundCpIdx)
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

func PostAppLogin(c echo.Context, reqAppLoginInfo *context.RequestAppLoginInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	appInfo := new(context.AppInfo)

	// 1. 인증 서버 접근
	if appID, CompanyID, returnValue, err := model.GetDB().GetApplications(&reqAppLoginInfo.Account); err != nil || returnValue != 1 {
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.SetReturn(resultcode.Result_Auth_NotMatchAppAccount)
		}
		return c.JSON(http.StatusOK, resp)
	} else {
		appInfo.AppID = appID
		appInfo.CompanyID = CompanyID
	}

	// 2. Access, Refresh 토큰 생성
	if jwtInfoValue, err := auth.GetIAuth().MakeToken(context.LoginType(context.AppLogin), appInfo); err != nil {
		resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
		return c.JSON(http.StatusOK, resp)
	} else {
		resp.Value = jwtInfoValue
	}

	return c.JSON(http.StatusOK, resp)
}

func DelAppLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	if err := auth.GetIAuth().DeleteJwtInfo(context.LoginType(context.AppLogin), ctx.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_RedisError)
	}

	return c.JSON(http.StatusOK, resp)
}
