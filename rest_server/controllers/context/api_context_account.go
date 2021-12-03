package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type Account struct {
	AUID       int    `json:"au_id"`
	SocialID   string `json:"social_id" validate:"required"`
	SocialType int    `json:"social_type"`
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

func (o *Account) CheckValidate() *base.BaseResponse {
	if len(o.SocialID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialInfo)
	}
	return nil
}
