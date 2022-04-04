package app

import (
	"fmt"
	"sync"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/api_inno_log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/externalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/internalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
)

type ServerApp struct {
	base.BaseApp
	conf       *config.ServerConfig
	configFile string

	auth *auth.IAuth
}

func (o *ServerApp) Init(configFile string) (err error) {
	o.conf = config.GetInstance(configFile)
	base.AppendReturnCodeText(&resultcode.ResultCodeText)
	context.AppendRequestParameter()

	if err := o.NewDB(o.conf); err != nil {
		return err
	}

	o.InitLogServer(o.conf)

	if auth, err := auth.NewIAuth(&o.conf.Auth); err != nil {
		return err
	} else {
		o.auth = auth
	}

	return err
}

func (o *ServerApp) InitLogServer(conf *config.ServerConfig) {
	confLog := conf.InnoLog
	hostInfo := api_inno_log.HostInfo{
		IntHostUri: confLog.InternalpiDomain,
		ExtHostUri: confLog.ExternalpiDomain,
		IntVer:     confLog.InternalVer,
		ExtVer:     confLog.ExternalVer,
	}
	api_inno_log.NewServerInfo("", hostInfo)
}

func (o *ServerApp) CleanUp() {
	fmt.Println("CleanUp")
}

func (o *ServerApp) Run(wg *sync.WaitGroup) error {
	return nil
}

func (o *ServerApp) GetConfig() *baseconf.Config {
	return &o.conf.Config
}

func NewApp() (*ServerApp, error) {
	app := &ServerApp{}

	intAPI := internalapi.NewAPI()
	extAPI := externalapi.NewAPI()

	if err := app.BaseApp.Init(app, intAPI, extAPI); err != nil {
		return nil, err
	}

	return app, nil
}

func (o *ServerApp) NewDB(conf *config.ServerConfig) error {
	return model.InitDB(conf)
}
