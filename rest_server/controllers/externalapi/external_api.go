package externalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/labstack/echo"
)

type ExternalAPI struct {
	base.BaseController

	conf    *config.ServerConfig
	apiConf *baseconf.APIServer
	echo    *echo.Echo
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
		appInfo, accessUuid, err := auth.GetIAuth().VerifyAccessToken(author[0][7:])
		if err != nil {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		base.GetContext(c).(*context.InnoAuthContext).SetAppInfo(appInfo, accessUuid)
	}

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func (o *ExternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	// serving documents for RESTful APIs
	if o.conf.InnoAuth.APIDocs {
		e.Static("/docs", "docs/ext")
	}

	return nil
}

func (o *ExternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[1]
	return o.apiConf
}

func NewAPI() *ExternalAPI {
	return &ExternalAPI{}
}

func (o *ExternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *ExternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}
