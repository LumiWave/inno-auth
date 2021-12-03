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
	account := context.NewAccount()

	if err := c.Bind(account); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.PostMemberLogin(c, account)
}

// App 존재 여부 확인
func (o *ExternalAPI) GetMemberExists(c echo.Context) error {
	account := context.NewAccount()

	if err := c.Bind(account); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}

	return commonapi.GetMemberExists(c, account)
}
