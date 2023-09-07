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
	resAccountWeb, err := model.GetDB().AuthAccounts(reqAccountWeb)
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
