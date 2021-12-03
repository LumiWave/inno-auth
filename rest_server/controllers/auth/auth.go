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

func (o *IAuth) MakeToken(loginType context.LoginType, app *context.Application) (*context.JwtInfo, error) {
	jwtInfo := &context.JwtInfo{
		AccessUuid:  uuid.NewV4().String(),
		RefreshUuid: uuid.NewV4().String(),

		AtExpireDt: time.Now().Add(time.Minute * time.Duration(o.conf.AccessTokenExpiryPeriod)).UnixMilli(),
		RtExpireDt: time.Now().Add(time.Minute * time.Duration(o.conf.RefreshTokenExpiryPeriod)).UnixMilli(),
	}

	//create access token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = jwtInfo.AccessUuid
	atClaims["login_type"] = loginType
	atClaims["cp_id"] = app.CompanyID
	atClaims["app_id"] = app.AppID
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
	rtClaims["login_type"] = loginType
	rtClaims["cp_id"] = app.CompanyID
	rtClaims["app_id"] = app.AppID
	rtClaims["exp"] = jwtInfo.RtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}

	//redis save
	jwtInfo.RefreshToken = refreshToken
	if err := o.SetJwtInfo(jwtInfo, loginType, app); err != nil {
		return nil, err
	}

	return jwtInfo, nil
}

// jwt verify check
func (o *IAuth) VerifyAccessToken(accessToken string) (*context.Application, context.LoginType, string, error) {
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
		return nil, context.LoginType(context.NoneLogin), "", err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, context.LoginType(context.NoneLogin), "", errors.New("invalid access jwt")
	}

	accessUuid := fmt.Sprintf("%v", atClaims["access_uuid"])
	loginType := context.LoginType(int(atClaims["login_type"].(float64)))

	app, err := o.GetJwtInfo(loginType, accessUuid)
	if err != nil {
		return nil, context.LoginType(context.NoneLogin), "", err
	}

	return app, loginType, accessUuid, nil
}

func (o *IAuth) VerifyRefreshToken(refreshToken string) (*context.Application, context.LoginType, string, error) {
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
		return nil, context.LoginType(context.NoneLogin), "", err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return nil, context.LoginType(context.NoneLogin), "", errors.New("invalid refresh jwt")
	}

	refreshUuid := fmt.Sprintf("%v", atClaims["refresh_uuid"])
	loginType := context.LoginType(int(atClaims["login_type"].(float64)))

	app, err := o.GetJwtInfo(loginType, refreshUuid)
	if err != nil {
		return nil, context.LoginType(context.NoneLogin), "", err
	}

	return app, loginType, refreshUuid, nil
}

// redis jwt info set
func (o *IAuth) SetJwtInfo(tokenInfo *context.JwtInfo, loginType context.LoginType, app *context.Application) error {
	return model.GetDB().SetJwtInfo(tokenInfo, loginType, app)
}

func (o *IAuth) GetJwtInfo(loginType context.LoginType, uuid string) (*context.Application, error) {
	return model.GetDB().GetJwtInfo(loginType, uuid)
}

func (o *IAuth) DeleteJwtInfo(loginType context.LoginType, uuid string) error {
	return model.GetDB().DeleteJwtInfo(loginType, uuid)
}
