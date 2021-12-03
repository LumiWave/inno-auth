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

func NewAccount() *Account {
	return new(Account)
}

func (o *Account) CheckValidate() *base.BaseResponse {
	if len(o.SocialID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyMemberSocialInfo)
	}
	return nil
}
