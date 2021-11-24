package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostMemberRegister(c echo.Context, memberInfo *context.MemberInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// Member Social Info 빈문자열 체크
	if err := memberInfo.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// Member Social Info 중복 체크
	if value, err := model.GetDB().SelectGetMemberInfoByASocialUID(memberInfo.Social.SocialUID); err == nil {
		if len(value.Social.SocialUID) > 0 {
			log.Error("PostMemberRegister exists social_uid", value.Social.SocialUID, " errorCode:", resultcode.Result_Auth_ExistsMemberSocialInfo)
			resp.SetReturn(resultcode.Result_Auth_ExistsMemberSocialInfo)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 테이블에 신규 row 생성
	if err := model.GetDB().InsertMember(memberInfo); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return c.JSON(http.StatusOK, resp)
}

func PostMemberLogin(c echo.Context, memberInfo *context.MemberInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 가입 정보 확인

	// 2. redis duplicate check
	// redis에 기존 정보가 있다면 기존에 발급된 토큰으로 응답한다.

	// 3. create Auth Token

	return c.JSON(http.StatusOK, resp)
}
