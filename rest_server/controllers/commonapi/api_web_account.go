package commonapi

import (
	"fmt"
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

// Web 계정 로그인/가입
func PostWebAccountLogin(c echo.Context, params *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance()

	for i := 21; i <= 40; i++ { // 여기 값(Social Key) 범위를 변경하면 됩니다.
		userStr := fmt.Sprintf("%021d", i)
		payload := &context.Payload{
			LoginType: context.WebAccountLogin,
			InnoUID: inno.AESEncrypt(inno.MakeInnoID(userStr, params.SocialType),
				[]byte(conf.Secret.Key),
				[]byte(conf.Secret.Iv)),
		}
		reqAccountWeb := &context.ReqAccountWeb{
			InnoUID:    payload.InnoUID,
			SocialID:   userStr,
			SocialType: auth.SocialType_Google,
			EA:         "",
		}

		// 2. 웹 로그인/가입
		resAccountWeb, needWallets, err := model.GetDB().AuthAccounts(reqAccountWeb)
		if err != nil {
			log.Error(err)
			resp.SetReturn(resultcode.Result_DBError)
			return c.JSON(http.StatusOK, resp)
		} else {
			payload.AUID = resAccountWeb.AUID
			resAccountWeb.InnoUID = payload.InnoUID
		}

		if len(needWallets) > 0 {
			walletInfo, err := inner.TokenAddressNew(needWallets, payload.InnoUID)
			if err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
				return c.JSON(http.StatusOK, resp)
			}

			// 3-2. [DB] ETH, MATIC 지갑 생성 프로시저 호출
			if err := model.GetDB().AddAccountBaseCoins(resAccountWeb.AUID, walletInfo); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Base_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}

			for _, needWallet := range needWallets {
				baseCoin, exists := model.GetDB().DBMeta.BaseCoins[needWallet.BaseCoinID]
				if !exists {
					log.Errorf("needWallet.BaseCoinID is not exists : %v", needWallet.BaseCoinID)
					resp.SetReturn(resultcode.Result_NeedWallet_BaseCoins_Error)
					return c.JSON(http.StatusOK, resp)
				}

				// 3-3. [DB] ONIT, ETH, MATIC 사용자 코인 등록
				if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, baseCoin.IDList); err != nil {
					log.Errorf("%v", err)
					resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
					return c.JSON(http.StatusOK, resp)
				}
			}
		}

	}

	return c.JSON(http.StatusOK, resp)
}

// Web 계정 로그아웃
func DelWebAccountLogout(c echo.Context) error {
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	resp := new(base.BaseResponse)
	resp.Success()

	// Check if the token has expired
	if _, err := auth.GetIAuth().GetJwtInfoByInnoUID(ctx.Payload.LoginType, context.AccessT, ctx.Payload.InnoUID); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		// Delete the innoUID in Redis.
		if err := auth.GetIAuth().DeleteInnoUIDRedis(ctx.Payload.LoginType, context.AccessT, ctx.Payload.InnoUID); err != nil {
			resp.SetReturn(resultcode.Result_RedisError)
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// Web 계정 로그인 정보 확인
func PostWebAccountInfo(c echo.Context, params *context.ReqAccountInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if jwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(context.WebAccountLogin, context.AccessT, params.InnoUID); err != nil {
		resp.SetReturn(resultcode.Result_Auth_ExpiredJwt)
	} else {
		respWebAccountInfo := &context.ResWebAccountInfo{
			JwtInfo:    *jwtInfo,
			InnoUID:    params.InnoUID,
			AUID:       params.AUID,
			SocialType: params.SocialType,
		}
		resp.Value = *respWebAccountInfo
	}

	return c.JSON(http.StatusOK, resp)
}
