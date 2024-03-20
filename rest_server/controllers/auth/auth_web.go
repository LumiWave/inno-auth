package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/hkdf"
)

func (o *IAuth) MakeWebToken(payload *context.Payload) (*context.JwtInfo, error) {
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
	if err := o.SetJwtInfoByInnoUID(jwtInfo, payload); err != nil {
		return nil, err
	}

	return jwtInfo, err
}

func (o *IAuth) DeleteInnoUIDRedis(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	// Redis에 AccessToken 정보 삭제
	if err := o.DeleteJwtInfoByInnoUID(loginType, context.AccessT, innoUID); err != nil {
		return err
	}

	// Redis에 RefreshToken 정보 삭제
	if err := o.DeleteJwtInfoByInnoUID(loginType, context.RefreshT, innoUID); err != nil {
		return err
	}
	return nil
}

func (o *IAuth) WebTokenRenew(payload *context.Payload) (*context.JwtInfo, int) {
	if jwtInfo, err := o.GetJwtInfoByInnoUID(payload.LoginType, context.RefreshT, payload.InnoUID); err != nil {
		return nil, resultcode.Result_Auth_ExpiredJwt
	} else {
		// 1. 기존 로그인 정보 (AccessToken, RefreshToken) 삭제
		if err := o.DeleteInnoUIDRedis(payload.LoginType, context.RefreshT, payload.InnoUID); err != nil {
			return nil, resultcode.Result_RedisError
		}
		// accesstoken payload에는 zklogin 관련정보가 없기 때문에 redis에서 가져와서 load 한다.
		payload.ZkLogin = jwtInfo.ZkLogin
		// 2. Web 토큰 재발급
		if newJwtInfo, err := o.MakeWebToken(payload); err != nil {
			return nil, resultcode.Result_Auth_MakeTokenError
		} else {
			return newJwtInfo, 0
		}
	}
}

// set redis jwt info
func (o *IAuth) SetJwtInfoByInnoUID(tokenInfo *context.JwtInfo, payload *context.Payload) error {
	return model.GetDB().SetJwtInfoByInnoUID(tokenInfo, payload)
}

// get redis jwt info
func (o *IAuth) GetJwtInfoByInnoUID(loginType context.LoginType, tokenType context.TokenType, innoUID string) (*context.JwtInfo, error) {
	return model.GetDB().GetJwtInfoByInnoUID(loginType, tokenType, innoUID)
}

// delete redis jwt info
func (o *IAuth) DeleteJwtInfoByInnoUID(loginType context.LoginType, tokenType context.TokenType, innoUID string) error {
	return model.GetDB().DeleteJwtInfoByInnoUID(loginType, tokenType, innoUID)
}

func (o *IAuth) MakeSalt(idToken string) (string, error) {
	// 토큰 파싱 (여기서는 서명 검증을 생략하고 클레임만 추출합니다)
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid claims")
	}

	for key, val := range claims {
		fmt.Printf("%s: %v\n", key, val)
	}

	aud := ""
	sub := ""
	if getAud, ok := claims["aud"]; !ok {
		log.Errorf("not exist aud id_token:%v", idToken)
	} else {
		aud = getAud.(string)
	}

	if getsub, ok := claims["sub"]; !ok {
		log.Errorf("not exist sub id_token : %v", idToken)
	} else {
		sub = getsub.(string)
	}

	return makeHKDFsha256(o.conf.AccessSecretKey, aud, sub)
}

func makeHKDFsha256(seed, salt, info string) (string, error) {
	ikm := []byte(seed)
	saltByte := []byte(salt)
	infoByte := []byte(info)

	// HKDF 객체 생성
	hkdf := hkdf.New(sha256.New, ikm, saltByte, infoByte)

	// 파생된 키를 저장할 슬라이스 생성 (여기서는 SHA-256 해시 크기와 동일하게 16 바이트로 설정)
	derivedKey := make([]byte, 16)

	// HKDF로부터 키를 읽어 derivedKey에 저장
	if _, err := io.ReadFull(hkdf, derivedKey); err != nil {
		return "", err
	}

	return new(big.Int).SetBytes(derivedKey).String(), nil
}
