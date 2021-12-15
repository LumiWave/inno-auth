package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) GetSocialList(c echo.Context) error {
	return commonapi.GetSocialList(c)
}
