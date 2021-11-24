package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

// Member 신규 가입
func (o *InternalAPI) PostMemberRegister(c echo.Context) error {
	memberInfo := context.NewMemberInfo()
	if err := c.Bind(memberInfo); err != nil {
		log.Error(err)
		return base.BaseJSONInternalServerError(c, err)
	}
	context.MakeDt(&memberInfo.CreateDt)

	return commonapi.PostMemberRegister(c, memberInfo)
}
