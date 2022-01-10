package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type IAuth struct {
	conf *config.ApiAuth

	SocialAuths map[int]SocialAuth
}

func NewIAuth(conf *config.ApiAuth) (*IAuth, error) {
	if gAuth == nil {
		gAuth = new(IAuth)
		gAuth.conf = conf
		gAuth.SocialAuths = make(map[int]SocialAuth)
		MakeSocialAuths(gAuth)
	}

	return gAuth, nil
}

func GetIAuth() *IAuth {
	return gAuth
}

func (o *IAuth) MakeToken(payload *context.Payload) (*context.JwtInfo, error) {
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

	//redis save
	jwtInfo.RefreshToken = refreshToken
	if err := o.SetJwtInfo(jwtInfo, payload); err != nil {
		return nil, err
	}

	return jwtInfo, err
}

// jwt verify check
func (o *IAuth) VerifyAccessToken(accessToken string) (*context.Payload, error) {
	atClaims := jwt.MapClaims{}
	jwtData, err := jwt.ParseWithClaims(accessToken, atClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("")
			}
			return []byte(o.conf.AccessSecretKey), nil
		})
	if err != nil {
		//exp가 만료되면 여기로 에러 리턴됨
		return nil, err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, errors.New("invalid access jwt")
	}
	payload := &context.Payload{
		CompanyID: int(atClaims["company_id"].(float64)),
		AppID:     int(atClaims["app_id"].(float64)),
		LoginType: context.LoginType(int(atClaims["login_type"].(float64))),
		Uuid:      fmt.Sprintf("%v", atClaims["access_uuid"]),
	}

	return payload, nil
}

func (o *IAuth) VerifyRefreshToken(refreshToken string) (*context.Payload, error) {
	atClaims := jwt.MapClaims{}
	jwtData, err := jwt.ParseWithClaims(refreshToken, atClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("")
			}
			return []byte(o.conf.RefreshSecretKey), nil
		})
	if err != nil {
		//exp가 만료되면 여기로 에러 리턴됨
		return nil, err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, errors.New("invalid refresh jwt")
	}

	payload := &context.Payload{
		CompanyID: int(atClaims["company_id"].(float64)),
		AppID:     int(atClaims["app_id"].(float64)),
		LoginType: context.LoginType(int(atClaims["login_type"].(float64))),
		Uuid:      fmt.Sprintf("%v", atClaims["refresh_uuid"]),
	}

	return payload, nil
}

// redis jwt info set
func (o *IAuth) SetJwtInfo(tokenInfo *context.JwtInfo, payload *context.Payload) error {
	return model.GetDB().SetJwtInfo(tokenInfo, payload)
}

func (o *IAuth) GetJwtInfo(payload *context.Payload) (*context.JwtInfo, error) {
	return model.GetDB().GetJwtInfo(payload)
}

func (o *IAuth) DeleteJwtInfo(payload *context.Payload) error {
	return model.GetDB().DeleteJwtInfo(payload)
}
