package commonapi

import (
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/inno_log"
	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/ip"
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
func PostWebAccountLogin(c echo.Context, params *context.AccountWeb, isExt bool) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance()

	if isExt {
		if !auth.CheckValidateExternal(params.SocialType) {
			resp.SetReturn(resultcode.Result_Auth_InvalidSocial_Type)
			return c.JSON(http.StatusOK, resp)
		}
	} else {
		if !auth.CheckValidateInternal(params.SocialType) {
			resp.SetReturn(resultcode.Result_Auth_InvalidSocial_Type)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 1. 소셜 정보 검증
	userID, ea, err := auth.GetIAuth().SocialAuths[params.SocialType].VerifySocialKey(params.SocialKey)
	if err != nil || len(userID) == 0 {
		log.Errorf("VerifySocialKey(SocialType:%v), userID(%v), ea(%v), err(%v)", params.SocialType, userID, ea, err)
		resp.SetReturn(resultcode.Result_Auth_VerifySocial_Key)
		return c.JSON(http.StatusOK, resp)
	}

	payload := &context.Payload{
		LoginType:  context.WebAccountLogin,
		SocialType: params.SocialType,
		InnoUID: inno.AESEncrypt(inno.MakeInnoID(userID, params.SocialType),
			[]byte(conf.Secret.Key),
			[]byte(conf.Secret.Iv)),
	}

	// 1-1. InnoUID 생성 에러 오류
	if len(payload.InnoUID) == 0 {
		log.Errorf("MakeInnoUID is empty : InnoUID(%v)", payload.InnoUID)
		resp.SetReturn(resultcode.Result_Auth_Invalid_InnoUID)
		return c.JSON(http.StatusOK, resp)
	}

	reqAccountWeb := &context.ReqAccountWeb{
		InnoUID:    payload.InnoUID,
		SocialID:   userID,
		SocialType: params.SocialType,
		EA:         inno.AESEncrypt(ea, []byte(conf.Secret.Key), []byte(conf.Secret.Iv)),
	}

	// 2. 웹 로그인/가입
	resAccountWeb, needWallets, err := model.GetDB().AuthAccounts(reqAccountWeb)
	if err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
		return c.JSON(http.StatusOK, resp)
	} else {
		payload.AUID = resAccountWeb.AUID
		resAccountWeb.InnoUID = payload.InnoUID
		resAccountWeb.SocialType = params.SocialType
	}

	// 3. [DB] 사용자 로그 등록
	// 3-1. 유저 IP로 국가 정보 판단
	countryCode, _ := ip.GetCountryByIp(params.IP, conf.AccessCountry.LocationFilePath)

	// 3-2. Log 전송
	inner.PostAccountAuthLog(&inno_log.AccountAuthLog{
		LogDt:       time.Now().Format("2006-01-02 15:04:05.000"),
		LogID:       context.AccountAuthLog_Auth,
		AUID:        resAccountWeb.AUID,
		InnoUID:     resAccountWeb.InnoUID,
		SocialID:    userID,
		SocialType:  params.SocialType,
		CountryCode: countryCode,
	}, resAccountWeb.IsJoined)

	// 4. ETH, MATIC 메인 지갑이 없는 유저는 지갑을 생성
	// 4-1. [token-manager] ETH, MATIC 지갑 생성
	if len(needWallets) > 0 {
		walletInfo, err := inner.TokenAddressNew(needWallets, payload.InnoUID)
		if err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
			return c.JSON(http.StatusOK, resp)
		}

		// 4-2. [DB] ETH, MATIC 지갑 생성 프로시저 호출
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

			// 4-3. [DB] ONIT, ETH, MATIC 사용자 코인 등록
			if len(baseCoin.IDList) != 0 {
				if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, baseCoin.IDList); err != nil {
					log.Errorf("%v", err)
					resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
					return c.JSON(http.StatusOK, resp)
				}
			} else {
				log.Warnf("new coinID list is null [basecoin:%v, symbol:%v]", baseCoin.BaseCoinID, baseCoin.BaseCoinSymbol)
			}
		}
	} else {
		// 지갑 생성은 필요없지만 사용자 코인이 등록이 필요한 유저들은 등록 처리
		// 4-1. [DB] 등록된 사용자 코인 조회
		if usedCoinMap, err := model.GetDB().GetListAccountCoins(payload.AUID); err != nil {
			resp.SetReturn(resultcode.Result_Get_List_AccountCoins_Scan_Error)
			return c.JSON(http.StatusOK, resp)
		} else {
			// 4-2. [config] 사용자 코인 등록이 되어야할 리스트 구성 (ETH + MATIC)
			idList := append(config.GetInstance().EthToken.IDList, config.GetInstance().MaticToken.IDList...)

			// 4-3. 사용자 코인 등록 리스트와 등록된 사용자 코인을 비교해서 누락된 IDList를 구성
			var needIDList []int64
			for _, id := range idList {
				_, ok := usedCoinMap[id]
				if !ok {
					needIDList = append(needIDList, id)
				}
			}

			// 4-4. [DB] 누락된 IDList 추가 사용자 코인 등록
			if len(needIDList) > 0 {
				if err := model.GetDB().AddAccountCoins(resAccountWeb.AUID, needIDList); err != nil {
					log.Errorf("%v", err)
					resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
					return c.JSON(http.StatusOK, resp)
				}
			}
		}
	}

	// 5. Access, Refresh 토큰 생성
	//5-1. 기존에 발급된 토큰이 있는지 확인
	if oldJwtInfo, err := auth.GetIAuth().GetJwtInfoByInnoUID(payload.LoginType, context.AccessT, payload.InnoUID); err != nil || oldJwtInfo == nil {
		// 5-2. 기존에 발급된 토큰이 없다면 토큰을 발급한다. (Redis 확인)
		if jwtInfoValue, err := auth.GetIAuth().MakeWebToken(payload); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			return c.JSON(http.StatusOK, resp)
		} else {
			// 5-3. 새로 발급된 토큰으로 응답
			resAccountWeb.JwtInfo = *jwtInfoValue
		}
	} else {
		// 5-2. 기존 발급된 토큰으로 응답
		resAccountWeb.JwtInfo = *oldJwtInfo
	}

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
