package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/auth"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAppLogin(c echo.Context, access *context.Access) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 0. 정검중 체크
	if status, err := model.GetDB().GetCacheStatus(); err != nil {
		log.Errorf("system check!")
		resp.SetReturn(resultcode.Result_SystemCheck)
		return c.JSON(http.StatusOK, resp)
	} else {
		if status.IsMaintenance != 0 {
			resp.SetReturn(resultcode.Result_SystemCheck)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 1. 인증 서버 접근
	if payload, returnValue, err := model.GetDB().GetApplications(access); err != nil || returnValue != 1 {
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.SetReturn(resultcode.Result_Auth_NotMatchAppAccount)
		}
	} else {
		if !payload.IsEnabled {
			resp.SetReturn(resultcode.Result_Auth_DeactivatedAccount)
		} else {
			// 2. Access, Refresh 토큰 생성
			if jwtInfoValue, err := auth.GetIAuth().MakeAppToken(payload); err != nil {
				log.Errorf("%v", err)
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

	// 1. Access Token이 만료되었는지 확인
	if jwtInfo, err := auth.GetIAuth().GetJwtInfoByUUID(ctx.Payload.LoginType, context.AccessT, ctx.Payload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// 2. Redis에 Uuid를 조회해서 삭제 (로그아웃 처리)
		if err := auth.GetIAuth().DeleteUuidRedis(jwtInfo, ctx.Payload.LoginType, context.AccessT, ctx.Payload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}
