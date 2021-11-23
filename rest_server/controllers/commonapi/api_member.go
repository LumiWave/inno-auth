package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

func PostMemberLogin(c echo.Context, params *context.MemberInfo) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 가입 정보 확인

	// 2. redis duplicate check
	// redis에 기존 정보가 있다면 기존에 발급된 토큰으로 응답한다.

	// 3. create Auth Token

	return c.JSON(http.StatusOK, resp)
}
