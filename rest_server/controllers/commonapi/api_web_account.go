package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostWebAccountLogin(c echo.Context, accountWeb *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 소셜 정보 검증
	userID, email, err := auth.GetIAuth().SocialAuths[auth.SocialType_Google].VerifySocialKey(accountWeb.SocialKey)
	if err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_Auth_VerifySocial_Key)
		return c.JSON(http.StatusOK, resp)
	}
	innoUID := inno.AESEncrypt(inno.MakeInnoID(userID, email),
		[]byte(config.GetInstance().Secret.Key),
		[]byte(config.GetInstance().Secret.Iv))

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    innoUID,
		SocialID:   userID,
		SocialType: auth.SocialType_Google,
	}

	// 2. 웹 로그인/가입
	resAccountWeb, payload, err := model.GetDB().AuthWebAccounts(reqAccountWeb)
	if err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	}

	// 3. Access, Refresh 토큰 생성
	if jwtInfoValue, err := auth.GetIAuth().MakeToken(payload); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
		return c.JSON(http.StatusOK, resp)
	} else {
		resAccountWeb.JwtInfo = *jwtInfoValue
	}

	resp.Value = *resAccountWeb

	return c.JSON(http.StatusOK, resp)
}

func DelWebAccountLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	payload := new(context.Payload)
	{
		payload.LoginType = ctx.Payload.LoginType
		payload.Uuid = ctx.Payload.Uuid
	}

	if err := auth.GetIAuth().DeleteJwtInfo(payload); err != nil {
		resp.SetReturn(resultcode.Result_RedisError)
	}

	return c.JSON(http.StatusOK, resp)
}
