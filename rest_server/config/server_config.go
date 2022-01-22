package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
)

var once sync.Once
var currentConfig *ServerConfig

type InnoAuthServer struct {
	ApplicationName string `json:"application_name" yaml:"application_name"`
	APIDocs         bool   `json:"api_docs" yaml:"api_docs"`
}

type ApiAuth struct {
	AuthEnable                  bool   `yaml:"auth_enable"`
	AccessSecretKey             string `yaml:"access_secret_key"`
	RefreshSecretKey            string `yaml:"refresh_secret_key"`
	AppAccessTokenExpiryPeriod  int64  `yaml:"app_access_token_expiry_period"`
	AppRefreshTokenExpiryPeriod int64  `yaml:"app_refresh_token_expiry_period"`
	WebAccessTokenExpiryPeriod  int64  `yaml:"web_access_token_expiry_period"`
	WebRefreshTokenExpiryPeriod int64  `yaml:"web_refresh_token_expiry_period"`
	SignExpiryPeriod            int64  `yaml:"sign_expiry_period"`
	AesKey                      string `yaml:"aes_key"`
}

type TokenManagerServer struct {
	Uri string `yaml:"uri"`
}

type PointManagerServer struct {
	Uri string `yaml:"uri"`
}
type SecretInfo struct {
	Key string `yaml:"key"`
	Iv  string `yaml:"iv"`
}

type ServerConfig struct {
	baseconf.Config `yaml:",inline"`

	InnoAuth           InnoAuthServer  `yaml:"inno_auth_server"`
	MssqlDBAccountAll  baseconf.DBAuth `yaml:"mssql_db_account"`
	MssqlDBAccountRead baseconf.DBAuth `yaml:"mssql_db_account_read"`

	Auth         ApiAuth            `yaml:"api_auth"`
	TokenManager TokenManagerServer `yaml:"token_manager"`
	PointManager PointManagerServer `yaml:"point_manager"`
	Secret       SecretInfo         `yaml:"secret"`
}

func GetInstance(filepath ...string) *ServerConfig {
	once.Do(func() {
		if len(filepath) <= 0 {
			panic(baseconf.ErrInitConfigFailed)
		}
		currentConfig = &ServerConfig{}
		if err := baseconf.Load(filepath[0], currentConfig); err != nil {
			currentConfig = nil
		} else {
			if os.Getenv("ASPNETCORE_PORT") != "" {
				port, _ := strconv.ParseInt(os.Getenv("ASPNETCORE_PORT"), 10, 32)
				currentConfig.APIServers[0].Port = int(port)
				currentConfig.APIServers[1].Port = int(port)
				fmt.Println(port)
			}
		}
	})

	return currentConfig
}
