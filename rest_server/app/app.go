package app

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
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
	auth := conf.MssqlDBAuth
	port, err := strconv.ParseInt(auth.Port, 10, 32)
	if err != nil {
		log.Errorf("db port error : %v", port)
		return err
	}
	mssqlDB, err := basedb.GetMssql("D-INNO-ACCOUNT01", " ", auth.ID, auth.Password, auth.Host, int(port))
	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v, %v, %v",
			err, auth.Host, auth.ID, auth.Password, auth.Database, auth.PoolSize, auth.IdleSize)
		return err
	}

	gCache := basedb.GetCache(&conf.Cache)

	model.SetDB(mssqlDB, gCache)

	return nil
}
