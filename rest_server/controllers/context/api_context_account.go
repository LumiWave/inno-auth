package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Account struct {
	AUID       int    `json:"au_id"`
	SocialID   string `json:"social_id" validate:"required"`
	SocialType int    `json:"social_type" validate:"required"`
}

type RequestAccountAuth struct {
	Account     Account `json:"account" validate:"required"`
	AccessToken string  `json:"access_token" validate:"required"`
}

type AccountCoin struct {
	AUID          int    `json:"au_id"`
	CoinID        int    `json:"coin_id"`
	WalletAddress string `json:"wallet_address"`
	Quantity      string `json:"quantity"`
}

func NewAccount() *Account {
	return new(Account)
}

func (o *RequestAccountAuth) CheckValidate() *base.BaseResponse {
	if len(o.Account.SocialID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialID)
	} else if o.Account.SocialType == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialType)
	} else if len(o.AccessToken) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialType)
	}
	return nil
}
