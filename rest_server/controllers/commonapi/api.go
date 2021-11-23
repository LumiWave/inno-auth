package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"

	"github.com/labstack/echo"
)

func GetHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func GetVersion(c echo.Context, maxVersion string) error {
	resp := base.BaseResponse{}

	resp.Value = map[string]interface{}{"version": maxVersion,
		"revision": base.AppVersion}
	resp.Success()

	return c.JSON(http.StatusOK, resp)
}
