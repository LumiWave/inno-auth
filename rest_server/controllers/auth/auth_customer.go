package auth

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func (o *IAuth) MakeCustomerToken(payload *context.CustomerPayload) (*context.JwtInfo, error) {
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
	atClaims["account_id"] = payload.AccountID
	atClaims["customer_id"] = payload.CustomerID
	atClaims["exp"] = jwtInfo.AtExpireDt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(o.conf.AccessSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.AccessToken = accessToken

	//create refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["access_uuid"] = jwtInfo.AccessUuid
	rtClaims["login_type"] = payload.LoginType
	rtClaims["account_id"] = payload.AccountID
	rtClaims["customer_id"] = payload.CustomerID
	rtClaims["exp"] = jwtInfo.AtExpireDt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(o.conf.RefreshSecretKey))
	if err != nil {
		return nil, err
	}
	jwtInfo.RefreshToken = refreshToken

	//redis save
	if err := o.SetJwtInfoByCustomerUUID(jwtInfo, payload); err != nil {
		return nil, err
	}

	return jwtInfo, err
}

// set redis jwt info
func (o *IAuth) SetJwtInfoByCustomerUUID(tokenInfo *context.JwtInfo, payload *context.CustomerPayload) error {
	return model.GetDB().SetJwtInfoByCustomerUUID(tokenInfo, payload)
}

// get redis jwt info
func (o *IAuth) GetJwtInfoByCustomerUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) (*context.JwtInfo, error) {
	return model.GetDB().GetJwtInfoByCustomerUUID(loginType, tokenType, uuid)
}

// delete redis jwt info
func (o *IAuth) DeleteJwtInfoByCustomerUUID(loginType context.LoginType, tokenType context.TokenType, uuid string) error {
	return model.GetDB().DeleteJwtInfoByCustomerUUID(loginType, tokenType, uuid)
}
