package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

type MemberInfo struct {
	Idx      int64      `json:"idx" query:"idx"`
	MemberID string     `json:"member_id"`
	AppIdx   int64      `json:"app_idx" validate:"required"`
	Social   SocialInfo `json:"social_info" validate:"required"`
	Token    JwtInfo
	CreateDt int64 `json:"create_dt"`
}

type ResponseMemberInfo struct {
	MemberID string     `json:"member_id"`
	AppIdx   int64      `json:"app_idx"`
	Social   SocialInfo `json:"social_info" validate:"required"`
}

func NewMemberInfo() *MemberInfo {
	return new(MemberInfo)
}

func (o *MemberInfo) CheckValidate() *base.BaseResponse {
	if len(o.Social.SocialUID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyMemberSocialInfo)
	}
	return nil
}
