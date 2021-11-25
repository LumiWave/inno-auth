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
}

func NewIAuth(conf *config.ApiAuth) (*IAuth, error) {
	if gAuth == nil {
		gAuth = new(IAuth)
		gAuth.conf = conf
	}
	return gAuth, nil
}

func GetIAuth() *IAuth {
	return gAuth
}

func (o *IAuth) MakeToken(appInfo *context.AppInfo) (*context.JwtInfo, error) {
	jwtInfo := &context.JwtInfo{
		AccessUuid:  uuid.NewV4().String(),
		RefreshUuid: uuid.NewV4().String(),

		AtExpireDt: time.Now().Add(time.Minute * time.Duration(o.conf.TokenExpiryPeriod)).Unix(),
		RtExpireDt: time.Now().Add(time.Minute * time.Duration(o.conf.TokenExpiryPeriod)).Unix(),
	}

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = jwtInfo.AccessUuid
	atClaims["idx"] = appInfo.Idx
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
	rtClaims["idx"] = appInfo.Idx
	rtClaims["exp"] = jwtInfo.RtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}

	//redis save
	jwtInfo.RefreshToken = refreshToken
	appInfo.Token = *jwtInfo
	if err := o.SetJwtInfo(jwtInfo, appInfo); err != nil {
		return nil, err
	}

	return jwtInfo, nil
}

// jwt verify check
func (o *IAuth) VerifyAccessToken(accessToken string) (*context.AppInfo, string, error) {
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
		return nil, "", err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, "", errors.New("invalid access jwt")
	}

	accessUuid := fmt.Sprintf("%v", jwtData.Claims.(jwt.MapClaims)["access_uuid"])

	appInfo, err := o.GetJwtInfo(accessUuid)
	if err != nil {
		return nil, "", err
	}

	return appInfo, accessUuid, nil
}

func (o *IAuth) VerifyRefreshToken(refreshToken string) (*context.AppInfo, string, error) {
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
		return nil, "", err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, "", errors.New("invalid refresh jwt")
	}

	refreshUuid := fmt.Sprintf("%v", jwtData.Claims.(jwt.MapClaims)["refresh_uuid"])

	appInfo, err := o.GetJwtInfo(refreshUuid)
	if err != nil {
		return nil, "", err
	}

	return appInfo, refreshUuid, nil
}

// redis jwt info set
func (o *IAuth) SetJwtInfo(tokenInfo *context.JwtInfo, appInfo *context.AppInfo) error {
	return model.GetDB().SetJwtInfo(tokenInfo, appInfo)
}

func (o *IAuth) GetJwtInfo(uuid string) (*context.AppInfo, error) {
	return model.GetDB().GetJwtInfo(uuid)
}

func (o *IAuth) DeleteJwtInfo(uuid string) error {
	return model.GetDB().DeleteJwtInfo(uuid)
}
