package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func GetMeta(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	respMetaInfo := &context.RespMetaInfo{
		Socials: model.GetDB().SocialList,
	}

	resp.Value = respMetaInfo

	return c.JSON(http.StatusOK, resp)
}
