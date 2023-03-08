package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

////////// web 로그인/가입 정보
type AccountWeb struct {
	SocialKey  string `json:"social_key" validate:"required"`
	SocialType int64  `json:"social_type" validate:"required"`
	IP         string `json:"ip" validate:"required"`
}

func NewAccountWeb() *AccountWeb {
	return new(AccountWeb)
}

func (o *AccountWeb) CheckValidate() *base.BaseResponse {
	if len(o.SocialKey) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialKey)
	}
	if o.SocialType == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyAccountSocialType)
	}
	return nil
}

////////////////////////////////////////////

////////// web 계정 로그인 정보

type ReqAccountInfo struct {
	AUID       int64  `json:"au_id"`
	InnoUID    string `json:"inno_uid"`
	SocialType int64  `json:"social_type"`
}

func (o *ReqAccountInfo) CheckValidate(ctx *InnoAuthContext) *base.BaseResponse {
	ctxValue := ctx.GetValue()
	if ctxValue != nil {
		o.AUID = ctxValue.AUID
		o.InnoUID = ctx.GetValue().InnoUID
		o.SocialType = ctx.GetValue().SocialType
	}
	return nil
}

type ResWebAccountInfo struct {
	JwtInfo
	InnoUID    string `json:"inno_uid" validate:"required"`
	AUID       int64  `json:"au_id" validate:"required"`
	SocialType int64  `json:"social_type" validate:"required"`
}

////////////////////////////////////////////

type ReqAccountWeb struct {
	InnoUID    string `json:"inno_uid" validate:"required"`
	SocialID   string `json:"social_id" validate:"required"`
	SocialType int64  `json:"social_type" validate:"required"`
	EA         string `json:"ea" validate:"required"`
}

type ResAccountWeb struct {
	JwtInfo
	InnoUID    string `json:"inno_uid" validate:"required"`
	IsJoined   bool   `json:"is_joined" validate:"required"`
	AUID       int64  `json:"au_id" validate:"required"`
	SocialType int64  `json:"social_type" validate:"required"`
}

////////////////////////////////////////////
