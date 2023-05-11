package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostCustomerLogin(c echo.Context, access *context.CustomerAccess) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// 1. 인증 서버 접근
	if payload, err := model.GetDB().GetCustomerAccountsByAccountID(access); err != nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		// 2. access, refresh 토큰 생성
		if jwtInfoValue, err := auth.GetIAuth().MakeCustomerToken(payload); err != nil {
			log.Errorf("%v", err)
			resp.SetReturn(resultcode.Result_Auth_MakeTokenError)
			return c.JSON(http.StatusOK, resp)
		} else {
			resp.Value = jwtInfoValue
		}
	}

	return c.JSON(http.StatusOK, resp)
}
