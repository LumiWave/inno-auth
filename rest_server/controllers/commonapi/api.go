package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/schedule"

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

func PreCheck(c echo.Context) base.PreCheckResponse {
	conf := config.GetInstance()
	if err := base.SetContext(c, &conf.Config, context.NewInnoAuthServerContext); err != nil {
		log.Error(err)
		return base.PreCheckResponse{
			IsSucceed: false,
		}
	}

	// auth token 검증
	if conf.Auth.AuthEnable {
		author, ok := c.Request().Header["Authorization"]
		if !ok {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		payload, err := auth.GetIAuth().VerifyAccessToken(author[0][7:])
		if err != nil {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		base.GetContext(c).(*context.InnoAuthContext).SetAuthContext(payload)
	}

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func GetNodeMetric(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	node := schedule.GetSystemMonitor().GetMetricInfo()
	resp.Value = node

	return c.JSON(http.StatusOK, resp)
}
