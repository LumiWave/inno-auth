package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

////////// web 계정 정보
type AccountWeb struct {
	InnoUID    string `json:"inno_uid" validate:"required"`
	SocialID   string `json:"social_id" validate:"required"`
	SocialType string `json:"social_type" validate:"required"`
}

func NewAccountWeb() *AccountWeb {
	return new(AccountWeb)
}

func (o *AccountWeb) CheckValidate() *base.BaseResponse {
	if len(o.InnoUID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyInnoID)
	}
	if len(o.SocialID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialID)
	}
	if len(o.SocialType) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialType)
	}
	return nil
}

////////////////////////////////////////////
