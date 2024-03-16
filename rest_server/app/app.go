package app

import (
	"fmt"
	"sync"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/baseapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/externalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/internalapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/log_server"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/point_server"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/sui_enoki_server"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/sui_prover"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/token_server"
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

	baseapi.InitHttpClient()

	if err := o.NewDB(o.conf); err != nil {
		return err
	}

	o.InitInnoLog(o.conf)
	o.InitPointManager(o.conf)
	o.InitTokenManager(o.conf)
	o.InitSuiProve(o.conf)
	o.InitSuiEnoki(o.conf)

	if auth, err := auth.NewIAuth(&o.conf.Auth); err != nil {
		return err
	} else {
		o.auth = auth
	}

	return err
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

func (o *ServerApp) InitInnoLog(conf *config.ServerConfig) error {
	return log_server.InitInnoLog(conf)
}

func (o *ServerApp) InitPointManager(conf *config.ServerConfig) error {
	return point_server.InitPointManager(conf)
}

func (o *ServerApp) InitTokenManager(conf *config.ServerConfig) error {
	return token_server.InitTokenManager(conf)
}

func (o *ServerApp) InitSuiProve(conf *config.ServerConfig) error {
	return sui_prover.InitSuiProveManager(conf)
}

func (o *ServerApp) InitSuiEnoki(conf *config.ServerConfig) error {
	return sui_enoki_server.InitSuiEnokiManager(conf)
}
