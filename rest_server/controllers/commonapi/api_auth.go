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

func GetTokenVerify(c echo.Context) error {
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

	ctx := base.GetContext(c).(*context.InnoAuthContext)

	switch ctx.Payload.LoginType {
	case context.AppLogin:
		if _, err := auth.GetIAuth().GetJwtInfoByUUID(ctx.Payload.LoginType, context.AccessT, ctx.Payload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
		} else {
			resp.Value = ctx.Payload
		}
	case context.WebAccountLogin:
		if jwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(ctx.Payload.LoginType, context.AccessT, ctx.Payload.InnoUID); err != nil {
			resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
		} else {
			// accesstoken의 payload에는 zklogin 관련 정보는 없기 때문에 redis에서 load해서 응답해준다.
			ctx.Payload.ZkLogin = jwtInfo.ZkLogin
			resp.Value = ctx.Payload
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func PostTokenRenew(c echo.Context, refreshTokenRequest *context.RenewTokenRequest) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. Refresh 토큰 유효성 검증
	if loginType, rtClaims, err := auth.GetIAuth().VerifyRefreshToken(refreshTokenRequest.RefreshToken); err != nil {
		resp.SetReturn(resultcode.Result_Auth_InvalidJwt)
	} else {
		// 2. Payload를 생성
		payload := auth.GetIAuth().ParseClaimsToPayload(loginType, context.RefreshT, rtClaims)

		switch loginType {
		case context.AppLogin:
			// 3. App 토큰 재발급/갱신
			if newJwtInfo, resultCode := auth.GetIAuth().AppTokenRenew(payload); resultCode != 0 {
				resp.SetReturn(resultCode)
			} else {
				resp.Value = newJwtInfo
			}
		case context.WebAccountLogin:
			// 3. Web 토큰 재발급/갱신
			if newJwtInfo, resultCode := auth.GetIAuth().WebTokenRenew(payload); resultCode != 0 {
				resp.SetReturn(resultCode)
			} else {
				resp.Value = newJwtInfo
			}
		}

	}
	return c.JSON(http.StatusOK, resp)
}
