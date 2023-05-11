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

func PostCustomerLogin(c echo.Context, access *context.CustomerAccess) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 인증 서버 접근
	if payload, err := model.GetDB().GetCustomerAccountsByAccountID(access); err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		// 2. access, refresh 토큰 생성
		if jwtInfoValue, err := auth.GetIAuth().MakeCustomerToken(payload); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			return c.JSON(http.StatusOK, resp)
		} else {
			resp.Value = jwtInfoValue
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func DelCustomerLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. Access Token이 만료되었는지 확인
	if jwtInfo, err := auth.GetIAuth().GetJwtInfoByCustomerUUID(ctx.CustomerPayload.LoginType, context.AccessT, ctx.CustomerPayload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// 2. Redis에 Uuid를 조회해서 삭제 (로그아웃 처리)
		if err := auth.GetIAuth().DeleteUuidRedis(jwtInfo, ctx.CustomerPayload.LoginType, context.AccessT, ctx.CustomerPayload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func GetCustomerTokenVerify(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	ctx := base.GetContext(c).(*context.InnoAuthContext)

	if _, err := auth.GetIAuth().GetJwtInfoByUUID(ctx.CustomerPayload.LoginType, context.AccessT, ctx.CustomerPayload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		resp.Value = ctx.CustomerPayload
	}

	return c.JSON(http.StatusOK, resp)
}
