package auth

import (
	"time"

	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-auth/rest_server/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func (o *IAuth) MakeGameToken(payload *context.Payload) (*context.JwtInfo, error) {
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	jwtInfo := &context.JwtInfo{
		AccessUuid:  uuid.NewV4().String(),
		RefreshUuid: uuid.NewV4().String(),

		AtExpireDt: func() int64 {
			if payload.SocialType == SocialType_Inno {
				return time.Now().Add(time.Duration(365 * 24 * time.Hour)).UnixMilli()
			}
			return time.Now().Add(time.Duration(accessExpiryPeriod)).UnixMilli()
		}(),
		RtExpireDt: func() int64 {
			if payload.SocialType == SocialType_Inno {
				return time.Now().Add(time.Duration(365 * 24 * time.Hour)).UnixMilli()
			}
			return time.Now().Add(time.Duration(refreshExpiryPeriod)).UnixMilli()
		}(),
		ZkLogin: context.ZkLogin{
			IDToken:            payload.ZkLogin.IDToken,
			EphemeralPublicKey: payload.EphemeralPublicKey,
			Salt:               payload.Salt,
			Epoch:              payload.Epoch,
			Randomness:         payload.Randomness,
			Privatekey:         payload.Privatekey,
		},
	}

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = jwtInfo.AccessUuid
	atClaims["login_type"] = payload.LoginType
	atClaims["inno_uid"] = payload.InnoUID
	atClaims["au_id"] = payload.AUID
	atClaims["social_type"] = payload.SocialType
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
	rtClaims["inno_uid"] = payload.InnoUID
	rtClaims["au_id"] = payload.AUID
	rtClaims["social_type"] = payload.SocialType
	rtClaims["exp"] = jwtInfo.RtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.RefreshToken = refreshToken

	//redis save
	if err := o.SetJwtInfoByInnoUIDGame(jwtInfo, payload); err != nil {
		return nil, err
	}

	return jwtInfo, err
}

// zkloing 정보 제거하고 redis 에 저장 하지 않는 용도
func (o *IAuth) MakeGameTokenForApp(payload *context.Payload) (*context.JwtInfo, error) {
	// Select ExpiryPeriod (App or Web)
	accessExpiryPeriod, refreshExpiryPeriod := context.GetTokenExpiryperiod(payload.LoginType)

	jwtInfo := &context.JwtInfo{
		AccessUuid:  uuid.NewV4().String(),
		RefreshUuid: uuid.NewV4().String(),

		AtExpireDt: func() int64 {
			if payload.SocialType == SocialType_Inno {
				return time.Now().Add(time.Duration(365 * 24 * time.Hour)).UnixMilli()
			}
			return time.Now().Add(time.Duration(accessExpiryPeriod)).UnixMilli()
		}(),
		RtExpireDt: func() int64 {
			if payload.SocialType == SocialType_Inno {
				return time.Now().Add(time.Duration(365 * 24 * time.Hour)).UnixMilli()
			}
			return time.Now().Add(time.Duration(refreshExpiryPeriod)).UnixMilli()
		}(),
	}

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = jwtInfo.AccessUuid
	atClaims["login_type"] = payload.LoginType
	atClaims["inno_uid"] = payload.InnoUID
	atClaims["au_id"] = payload.AUID
	atClaims["social_type"] = payload.SocialType
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
	rtClaims["inno_uid"] = payload.InnoUID
	rtClaims["au_id"] = payload.AUID
	rtClaims["social_type"] = payload.SocialType
	rtClaims["exp"] = jwtInfo.RtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.RefreshToken = refreshToken

	return jwtInfo, err
}

func (o *IAuth) DeleteInnoUIDRedisGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	// Redis에 AccessToken 정보 삭제
	if err := o.DeleteJwtInfoByInnoUIDGame(loginType, context.AccessT, innoUID); err != nil {
		return err
	}

	// Redis에 RefreshToken 정보 삭제
	if err := o.DeleteJwtInfoByInnoUIDGame(loginType, context.RefreshT, innoUID); err != nil {
		return err
	}
	return nil
}

func (o *IAuth) GameTokenRenew(payload *context.Payload) (*context.JwtInfo, int) {
	if jwtInfo, err := o.GetJwtInfoByInnoUIDGame(payload.LoginType, context.RefreshT, payload.InnoUID); err != nil {
		return nil, resultcode.Result_Auth_ExpiredJwt
	} else {
		// 1. 기존 로그인 정보 (AccessToken, RefreshToken) 삭제
		if err := o.DeleteInnoUIDRedisGame(payload.LoginType, context.RefreshT, payload.InnoUID); err != nil {
			return nil, resultcode.Result_RedisError
		}
		// accesstoken payload에는 zklogin 관련정보가 없기 때문에 redis에서 가져와서 load 한다.
		payload.ZkLogin = jwtInfo.ZkLogin
		// 2. Web 토큰 재발급
		if newJwtInfo, err := o.MakeGameToken(payload); err != nil {
			return nil, resultcode.Result_Auth_MakeTokenError
		} else {
			return newJwtInfo, 0
		}
	}
}

// set redis jwt info
func (o *IAuth) SetJwtInfoByInnoUIDGame(tokenInfo *context.JwtInfo, payload *context.Payload) error {
	return model.GetDB().SetJwtInfoByInnoUIDGame(tokenInfo, payload)
}

// get redis jwt info
func (o *IAuth) GetJwtInfoByInnoUIDGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) (*context.JwtInfo, error) {
	return model.GetDB().GetJwtInfoByInnoUIDGame(loginType, tokenType, innoUID)
}

// delete redis jwt info
func (o *IAuth) DeleteJwtInfoByInnoUIDGame(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	return model.GetDB().DeleteJwtInfoByInnoUIDGame(loginType, tokenType, innoUID)
}
