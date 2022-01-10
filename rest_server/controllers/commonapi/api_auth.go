package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/labstack/echo"
)

func GetTokenVerify(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	ctx := base.GetContext(c).(*context.InnoAuthContext)

	if _, err := auth.GetIAuth().GetJwtInfo(ctx.Payload.LoginType, ctx.Payload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		resp.Value = ctx.Payload
	}

	return c.JSON(http.StatusOK, resp)
}

func PostTokenRenew(c echo.Context, refreshTokenRequest *context.RenewTokenRequest) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if payload, err := auth.GetIAuth().VerifyRefreshToken(refreshTokenRequest.RefreshToken); err != nil {
		resp.SetReturn(resultcode.Result_Auth_InvalidJwt)
	} else {
		// Get Jwt Info
		if _, err := auth.GetIAuth().GetJwtInfo(payload.LoginType, payload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
		} else {
			// Make Renew Token.
			if newJwtInfoValue, err := auth.GetIAuth().MakeToken(payload); err != nil {
				resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			} else {
				// Delete the uuid in Redis.
				if err := auth.GetIAuth().DeleteUuidRedis(payload.LoginType, payload.Uuid); err != nil {
					resp.SetReturn(resultcode.Result_RedisError)
				} else {
					resp.Value = newJwtInfoValue
				}
			}
		}
	}
	return c.JSON(http.StatusOK, resp)
}
