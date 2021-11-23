package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
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
	if value, err := model.GetDB().SelectApp(appInfo); err == nil {
		if len(value.AppName) > 0 {
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
