package config

import (
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
	AuthEnable               bool   `yaml:"auth_enable"`
	AccessSecretKey          string `yaml:"access_secret_key"`
	RefreshSecretKey         string `yaml:"refresh_secret_key"`
	AccessTokenExpiryPeriod  int64  `yaml:"access_token_expiry_period"`
	RefreshTokenExpiryPeriod int64  `yaml:"refresh_token_expiry_period"`
	SignExpiryPeriod         int64  `yaml:"sign_expiry_period"`
	AesKey                   string `yaml:"aes_key"`
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

	InnoAuth     InnoAuthServer     `yaml:"inno_auth_server"`
	MysqlDBAuth  baseconf.DBAuth    `yaml:"mysql_db_auth"`
	MssqlDBAuth  baseconf.DBAuth    `yaml:"mssql_db_auth"`
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
		}
	})

	return currentConfig
}
