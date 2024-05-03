package externalapi

import (
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/labstack/echo"
)

func (o *ExternalAPI) GetMeta(c echo.Context) error {
	return commonapi.GetMeta(c)
}
