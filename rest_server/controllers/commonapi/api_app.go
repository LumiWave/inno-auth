package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAppLogin(c echo.Context, access *context.Access) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 인증 서버 접근
	if payload, returnValue, err := model.GetDB().GetApplications(access); err != nil || returnValue != 1 {
		if err != nil {
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.SetReturn(resultcode.Result_Auth_NotMatchAppAccount)
		}
	} else {
		if !payload.IsEnabled {
			resp.SetReturn(resultcode.Result_Auth_DeactivatedAccount)
		} else {
			// 2. Access, Refresh 토큰 생성
			if jwtInfoValue, err := auth.GetIAuth().MakeToken(payload); err != nil {
				resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
				return c.JSON(http.StatusOK, resp)
			} else {
				resp.Value = jwtInfoValue
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func DelAppLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	payload := new(context.Payload)
	{
		payload.LoginType = context.AppLogin
		payload.Uuid = ctx.Payload.Uuid
	}

	if err := auth.GetIAuth().DeleteJwtInfo(payload); err != nil {
		resp.SetReturn(resultcode.Result_RedisError)
	}

	return c.JSON(http.StatusOK, resp)
}
