package auth

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func (o *IAuth) MakeAppToken(payload *context.Payload) (*context.JwtInfo, error) {
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	jwtInfo := &context.JwtInfo{
		AccessUuid:  uuid.NewV4().String(),
		RefreshUuid: uuid.NewV4().String(),

		AtExpireDt: time.Now().Add(time.Duration(accessExpiryPeriod)).UnixMilli(),
		RtExpireDt: time.Now().Add(time.Duration(refreshExpiryPeriod)).UnixMilli(),
	}
	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = jwtInfo.AccessUuid
	atClaims["login_type"] = payload.LoginType
	atClaims["company_id"] = payload.CompanyID
	atClaims["app_id"] = payload.AppID
	atClaims["exp"] = jwtInfo.AtExpireDt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(o.conf.AccessSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.AccessToken = accessToken

	//create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = jwtInfo.RefreshUuid
	rtClaims["login_type"] = payload.LoginType
	rtClaims["company_id"] = payload.CompanyID
	rtClaims["app_id"] = payload.AppID
	rtClaims["exp"] = jwtInfo.RtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.RefreshToken = refreshToken

	//redis save
	if err := o.SetJwtInfoByUUID(jwtInfo, payload); err != nil {
		return nil, err
	}

	return jwtInfo, err
}

func (o *IAuth) AppTokenRenew(payload *context.Payload) (*context.JwtInfo, int) {
	// Redis에서 Uuid로 jwtInfo 조회
	if jwtInfo, err := o.GetJwtInfoByUUID(payload.LoginType, context.RefreshT, payload.Uuid); err != nil {
		return nil, resultcode.Result_Auth_ExpiredJwt
	} else {
		// 기존 로그인 정보 (AccessToken, RefreshToken) 삭제
		if err := o.DeleteUuidRedis(jwtInfo, payload.LoginType, context.RefreshT, payload.Uuid); err != nil {
			return nil, resultcode.Result_RedisError
		}
		// App 토큰 재발급
		if newJwtInfo, err := o.MakeAppToken(payload); err != nil {
			return nil, resultcode.Result_Auth_MakeTokenError
		} else {
			return newJwtInfo, 0
		}
	}
}

func (o *IAuth) DeleteUuidRedis(jwtInfo *context.JwtInfo, loginType context.LoginType, tokenType context.TokenType, uuid string) error {
	// Redis에서 AccessToken 삭제
	if err := o.DeleteJwtInfoByUUID(loginType, context.AccessT, jwtInfo.AccessUuid); err != nil {
		return err
	}

	// Redis에서 RefreshToken 삭제
	if err := o.DeleteJwtInfoByUUID(loginType, context.RefreshT, jwtInfo.RefreshUuid); err != nil {
		return err
	}
	return nil
}

// set redis jwt info
func (o *IAuth) SetJwtInfoByUUID(tokenInfo *context.JwtInfo, payload *context.Payload) error {
	return model.GetDB().SetJwtInfoByUUID(tokenInfo, payload)
}

// get redis jwt info
func (o *IAuth) GetJwtInfoByUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) (*context.JwtInfo, error) {
	return model.GetDB().GetJwtInfoByUUID(loginType, tokenType, uuid)
}

// delete redis jwt info
func (o *IAuth) DeleteJwtInfoByUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) error {
	return model.GetDB().DeleteJwtInfoByUUID(loginType, tokenType, uuid)
}
