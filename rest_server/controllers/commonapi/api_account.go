package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAccountLogin(c echo.Context, reqAuthAccountApp *context.ReqAuthAccountApplication) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 인증 프로시저 호출 (신규 유저, 기존 유저를 체크)
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	if _, err := model.GetDB().AccountAuthApplication(reqAuthAccountApp, ctx.Payload); err != nil {
		log.Error(err)
	} else {

	}

	return c.JSON(http.StatusOK, resp)
}
