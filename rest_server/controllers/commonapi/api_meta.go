package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func GetSocialList(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	resp.Value = model.GetDB().SocialsS

	return c.JSON(http.StatusOK, resp)
}
