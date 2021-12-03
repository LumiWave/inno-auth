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
	resp.Value = ctx.Application()

	return c.JSON(http.StatusOK, resp)
}

func PostTokenRenew(c echo.Context, refreshTokenRequest *context.RenewTokenRequest) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if app, loginType, uuid, err := auth.GetIAuth().VerifyRefreshToken(refreshTokenRequest.RefreshToken); err != nil {
		resp.SetReturn(resultcode.Result_Auth_InvalidJwt)
	} else {
		// Make Renew Token.
		if jwtInfoValue, err := auth.GetIAuth().MakeToken(loginType, app); err != nil {
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
		} else {
			// Delete the uuid in Redis.
			auth.GetIAuth().DeleteJwtInfo(loginType, uuid)
			resp.Value = jwtInfoValue
		}
	}
	return c.JSON(http.StatusOK, resp)
}
