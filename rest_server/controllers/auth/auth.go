package auth

import (
	"errors"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/dgrijalva/jwt-go"
)

type IAuth struct {
	conf        *config.ApiAuth
	SocialAuths map[int64]SocialAuth
}

func NewIAuth(conf *config.ApiAuth) (*IAuth, error) {
	if gAuth == nil {
		gAuth = new(IAuth)
		gAuth.conf = conf
		gAuth.SocialAuths = make(map[int64]SocialAuth)
		MakeSocialAuths(gAuth)
	}

	return gAuth, nil
}

func GetIAuth() *IAuth {
	return gAuth
}

// jwt verify check
func (o *IAuth) VerifyAccessToken(accessToken string) (context.LoginType, jwt.MapClaims, error) {
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
		return context.NoneLogin, nil, err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return context.NoneLogin, nil, errors.New("invalid access jwt")
	}

	return context.LoginType(int(atClaims["login_type"].(float64))), atClaims, nil
}

func (o *IAuth) VerifyRefreshToken(refreshToken string) (context.LoginType, jwt.MapClaims, error) {
	rtClaims := jwt.MapClaims{}
	jwtData, err := jwt.ParseWithClaims(refreshToken, rtClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("")
			}
			return []byte(o.conf.RefreshSecretKey), nil
		})
	if err != nil {
		//exp가 만료되면 여기로 에러 리턴됨
		return context.NoneLogin, nil, err
	}

	if _, ok := jwtData.Claims.(jwt.MapClaims); !ok && !jwtData.Valid {
		return context.NoneLogin, nil, errors.New("invalid refresh jwt")
	}

	return context.LoginType(int(rtClaims["login_type"].(float64))), rtClaims, nil
}

func (o *IAuth) ParseClaimsToPayload(loginType context.LoginType, tokenType context.TokenType, claims jwt.MapClaims) *context.Payload {
	payload := new(context.Payload)

	var claimsType string
	if tokenType == context.AccessT {
		claimsType = "access_uuid"
	} else if tokenType == context.RefreshT {
		claimsType = "refresh_uuid"
	}

	switch loginType {
	case context.AppLogin:
		payload = &context.Payload{
			CompanyID: int64(claims["company_id"].(float64)),
			AppID:     int64(claims["app_id"].(float64)),
			LoginType: context.LoginType(int(claims["login_type"].(float64))),
			Uuid:      fmt.Sprintf("%v", claims[claimsType]),
		}
	case context.WebAccountLogin:
		payload = &context.Payload{
			LoginType:                  context.LoginType(int(claims["login_type"].(float64))),
			InnoUID:                    fmt.Sprintf("%v", claims["inno_uid"]),
			AUID:                       int64(claims["au_id"].(float64)),
			SocialType:                 int64(claims["social_type"].(float64)),
			Uuid:                       fmt.Sprintf("%v", claims[claimsType]),
			IDToken:                    fmt.Sprintf("%v", claims["id_token"]),
			ExtendedEphemeralPublicKey: fmt.Sprintf("%v", claims["extendedEphemeralPublicKey"]),
			EphemeralPublicKey:         fmt.Sprintf("%v", claims["ephemeralPublicKey"]),
			Salt:                       fmt.Sprintf("%v", claims["salt"]),
		}
	}
	return payload
}

func (o *IAuth) ParseClaimsToCustomerPayload(loginType context.LoginType, tokenType context.TokenType, claims jwt.MapClaims) *context.CustomerPayload {
	payload := new(context.CustomerPayload)

	var claimsType string
	if tokenType == context.AccessT {
		claimsType = "access_uuid"
	} else if tokenType == context.RefreshT {
		claimsType = "refresh_uuid"
	}

	payload = &context.CustomerPayload{
		AccountID:  int64(claims["account_id"].(float64)),
		CustomerID: int64(claims["customer_id"].(float64)),
		LoginType:  context.LoginType(int(claims["login_type"].(float64))),
		Uuid:       fmt.Sprintf("%v", claims[claimsType]),
	}

	return payload
}
