package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// 회원 로그인
func (o *ExternalAPI) PostMemberLogin(c echo.Context) error {
	params := context.NewMemberInfo()

	if err := c.Bind(params); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostMemberLogin(c, params)
}
