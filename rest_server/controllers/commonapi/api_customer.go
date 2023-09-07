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

	// 1. password md5 암호화
	access.AccessPW = auth.GetMd5Hash(access.AccessPW)

	// 2. 인증 서버 접근
	if payload, returnValue, err := model.GetDB().GetCustomerAccountsByAccountID(access); err != nil || returnValue != 1 {
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_DBError)
		} else {
			resp.SetReturn(resultcode.Result_Auth_Invalid_Customer_AccountID)
		}
	} else {
		// 1. access, refresh 토큰 생성
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

func PostCustomerTokenRenew(c echo.Context, refreshTokenRequest *context.RenewTokenRequest) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. Refresh 토큰 유효성 검증
	if loginType, rtClaims, err := auth.GetIAuth().VerifyRefreshToken(refreshTokenRequest.RefreshToken); err != nil {
		resp.SetReturn(resultcode.Result_Auth_InvalidJwt)
	} else {
		// 2. Payload를 생성
		payload := auth.GetIAuth().ParseClaimsToCustomerPayload(loginType, context.RefreshT, rtClaims)

		// 3. Customer 토큰 재발급/갱신
		if newJwtInfo, resultCode := auth.GetIAuth().CustomerTokenRenew(payload); resultCode != 0 {
			resp.SetReturn(resultCode)
		} else {
			resp.Value = newJwtInfo
		}

	}
	return c.JSON(http.StatusOK, resp)
}
