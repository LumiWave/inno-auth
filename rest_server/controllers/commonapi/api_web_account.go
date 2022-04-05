package commonapi

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

const (
	AccountAuthLog_NewAccount = 5
	AccountAuthLog_Account    = 6
)

// Web 계정 로그인/가입
func PostWebAccountLogin(c echo.Context, params *context.AccountWeb) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance()

	// 1. 소셜 정보 검증
	// userID, email, err := auth.GetIAuth().SocialAuths[params.SocialType].VerifySocialKey(params.SocialKey)
	// if err != nil || len(userID) == 0 || len(email) == 0 {
	// 	log.Errorf("%v", err)
	// 	resp.SetReturn(resultcode.Result_Auth_VerifySocial_Key)
	// 	return c.JSON(http.StatusOK, resp)
	// }

	userID := strconv.FormatInt(time.Now().UnixNano()+rand.Int63(), 10)
	//startTime := time.Now().UnixMilli()
	//_, _, err := auth.GetIAuth().SocialAuths[params.SocialType].VerifySocialKey(params.SocialKey)
	//log.Errorf("Google VerifySocialKey time %v", time.Now().UnixMilli()-startTime)

	//startTime1 := time.Now().UnixMilli()
	payload := &context.Payload{
		LoginType:  context.WebAccountLogin,
		SocialType: params.SocialType,
		InnoUID: inno.AESEncrypt(inno.MakeInnoID(userID, params.SocialType),
			[]byte(conf.Secret.Key),
			[]byte(conf.Secret.Iv)),
	}
	// endTime1 := time.Now().UnixMilli()
	// if endTime1-startTime1 >= 1000 {
	// 	log.Errorf("MakeInnoID time:%v", endTime1-startTime1)
	// }

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    payload.InnoUID,
		SocialID:   userID,
		SocialType: params.SocialType,
	}

	// 2. 웹 로그인/가입
	// startTime2 := time.Now().UnixMilli()
	resAccountWeb, err := model.GetDB().AuthAccounts(reqAccountWeb)
	if err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	}
	// endTime2 := time.Now().UnixMilli()
	// if endTime2-startTime2 >= 1000 {
	// 	log.Errorf("AuthAccounts time: %v", endTime2-startTime2)
	// }
	payload.AUID = resAccountWeb.AUID
	resAccountWeb.InnoUID = payload.InnoUID
	resAccountWeb.SocialType = params.SocialType

	// 4. ONIT 지갑이 없는 유저는 지갑을 생성
	if !resAccountWeb.ExistsMainWallet {
		// 3-1. [token-manager] ETH 지갑 생성
		var baseCoinList []context.CoinInfo
		for i, value := range conf.BaseCoin.SymbolList {
			baseCoinList = append(baseCoinList, context.CoinInfo{
				CoinID:     conf.BaseCoin.IDList[i],
				CoinSymbol: value,
			})
		}
		// startTime3 := time.Now().UnixMilli()
		walletInfo, err := inner.TokenAddressNew(baseCoinList, payload.InnoUID)
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
			return c.JSON(http.StatusOK, resp)
		}

		// endTime3 := time.Now().UnixMilli()
		// if endTime3-startTime3 >= 1000 {
		// 	log.Errorf("%v", endTime3-startTime3)
		// }

		// 3-2. [DB] ETH 지갑 생성 프로시저 호출
		// startTime4 := time.Now().UnixMilli()
		if err := model.GetDB().AddAccountBaseCoins(resAccountWeb.AUID, walletInfo); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Procedure_Add_Base_Account_Coins)
			return c.JSON(http.StatusOK, resp)
		}

		// endTime4 := time.Now().UnixMilli()
		// if endTime4-startTime4 >= 1000 {
		// 	log.Errorf("AddAccountBaseCoins time: %v", endTime4-startTime4)
		// }

		// startTime5 := time.Now().UnixMilli()
		// 3-3. [DB] ONIT 사용자 코인 등록
		if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, conf.ProjectToken.IDList); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
			return c.JSON(http.StatusOK, resp)
		}
		// endTime5 := time.Now().UnixMilli()
		// if endTime5-startTime5 >= 1000 {
		// 	log.Errorf("AddAccountCoins time: %v", endTime5-startTime5)
		// }
	}

	// 4. Access, Refresh 토큰 생성
	// startTime6 := time.Now().UnixMilli()
	//4-1. 기존에 발급된 토큰이 있는지 확인
	if oldJwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(payload.LoginType, context.AccessT, payload.InnoUID); err != nil || oldJwtInfo == nil {
		// 4-2. 기존에 발급된 토큰이 없다면 토큰을 발급한다. (Redis 확인)
		if jwtInfoValue, err := auth.GetIAuth().MakeWebToken(payload); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			return c.JSON(http.StatusOK, resp)
		} else {
			// 4-3. 새로 발급된 토큰으로 응답
			resAccountWeb.JwtInfo = *jwtInfoValue
		}
	} else {
		// 4-2. 기존 발급된 토큰으로 응답
		resAccountWeb.JwtInfo = *oldJwtInfo
	}
	// endTime6 := time.Now().UnixMilli()
	// if endTime6-startTime6 >= 1000 {
	// 	log.Errorf("MakeWebToken time: %v", endTime6-startTime6)
	// }

	resp.Value = *resAccountWeb

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
