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

func PostAppRegister(c echo.Context, params *context.AppInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// App 이름 빈문자열 체크
	if err := params.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// App 이름 중복 체크
	if value, err := model.GetDB().SelectApp(params); err == nil {
		if len(value.AppName) > 0 {
			log.Error("PostAppRegister exists app_name=", params.AppName, " errorCode:", resultcode.Result_Auth_ExistsAppName)
			resp.SetReturn(resultcode.Result_Auth_ExistsAppName)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// cp_idx 존재 여부 확인
	if value, err := model.GetDB().SelectGetCPInfo(params.CpIdx); err == nil {
		if len(value.CpName) == 0 {
			log.Error("PostCPRegister Not Exists cp_idx=", params.CpIdx, " errorCode:", resultcode.Result_Auth_NotFoundCpIdx)
			resp.SetReturn(resultcode.Result_Auth_NotFoundCpName)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 테이블에 신규 row 생성
	if err := model.GetDB().InsertApp(params); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return c.JSON(http.StatusOK, resp)
}
