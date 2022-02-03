package internalapi

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/labstack/echo"
)

type InternalAPI struct {
	base.BaseController

	conf    *config.ServerConfig
	apiConf *baseconf.APIServer
	echo    *echo.Echo
}

func (o *InternalAPI) Init(e *echo.Echo) error {
	o.echo = e
	o.BaseController.PreCheck = commonapi.PreCheck

	if err := o.MapRoutes(o, e, o.apiConf.Routes); err != nil {
		return err
	}

	return nil
}

func (o *InternalAPI) GetConfig() *baseconf.APIServer {
	o.conf = config.GetInstance()
	o.apiConf = &o.conf.APIServers[0]
	return o.apiConf
}

func NewAPI() *InternalAPI {
	return &InternalAPI{}
}

func (o *InternalAPI) GetHealthCheck(c echo.Context) error {
	return commonapi.GetHealthCheck(c)
}

func (o *InternalAPI) GetVersion(c echo.Context) error {
	return commonapi.GetVersion(c, o.BaseController.MaxVersion)
}

func (o *InternalAPI) GetNodeMetric(c echo.Context) error {
	return commonapi.GetNodeMetric(c)
}

func (o *InternalAPI) GetInnoUIDInfo(c echo.Context) error {
	req := new(context.ReqGetInnoUID)

	if err := c.Bind(req); err != nil {
		log.Errorf("%v", err)
		return base.BaseJSONInternalServerError(c, err)
	}
	return commonapi.GetInnoUIDInfo(c, req.InnoUID)
}
