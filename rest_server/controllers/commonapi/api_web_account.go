package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostWebAccountLogin(c echo.Context, accountWeb *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance()

	// 1. 소셜 정보 검증
	userID, email, err := auth.GetIAuth().SocialAuths[accountWeb.SocialType].VerifySocialKey(accountWeb.SocialKey)
	if err != nil || len(userID) == 0 || len(email) == 0 {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_Auth_VerifySocial_Key)
		return c.JSON(http.StatusOK, resp)
	}

	payload := &context.Payload{
		LoginType: context.WebAccountLogin,
		InnoUID: inno.AESEncrypt(inno.MakeInnoID(userID, email),
			[]byte(conf.Secret.Key),
			[]byte(conf.Secret.Iv)),
	}

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    payload.InnoUID,
		SocialID:   userID,
		SocialType: accountWeb.SocialType,
	}

	// 2. 웹 로그인/가입
	resAccountWeb, err := model.GetDB().AuthAccounts(reqAccountWeb)
	if err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	}
	payload.AUID = resAccountWeb.AUID
	resAccountWeb.InnoUID = payload.InnoUID

	// 3. ONIT 지갑이 없는 유저는 지갑을 생성
	if !resAccountWeb.ExistsMainWallet {
		// 3-1. token-manager에 새 지갑 주소 생성 요청
		coinList := []context.CoinInfo{{
			CoinID:   conf.ONIT.ID,
			CoinName: conf.ONIT.Symbol,
		}}

		addressList, err := inner.TokenAddressNew(coinList, payload.InnoUID)
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
			return c.JSON(http.StatusOK, resp)
		}

		// 3-2. [DB] 지갑 생성 프로시저 호출
		if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, addressList); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 4. Access, Refresh 토큰 생성
	if jwtInfoValue, err := auth.GetIAuth().MakeToken(payload); err != nil {
		log.Errorf("%v", err)
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

	// Check if the token has expired
	if _, err := auth.GetIAuth().GetJwtInfo(ctx.Payload.LoginType, ctx.Payload.Uuid); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// Delete the uuid in Redis.
		if err := auth.GetIAuth().DeleteUuidRedis(ctx.Payload.LoginType, ctx.Payload.Uuid); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}
