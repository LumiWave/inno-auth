package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	AuthEnable                       bool   `yaml:"auth_enable"`
	AccessSecretKey                  string `yaml:"access_secret_key"`
	RefreshSecretKey                 string `yaml:"refresh_secret_key"`
	CustomerAccessTokenExpiryPeriod  int64  `yaml:"customer_access_token_expiry_period"`
	CustomerRefreshTokenExpiryPeriod int64  `yaml:"customer_refresh_token_expiry_period"`
	AppAccessTokenExpiryPeriod       int64  `yaml:"app_access_token_expiry_period"`
	AppRefreshTokenExpiryPeriod      int64  `yaml:"app_refresh_token_expiry_period"`
	WebAccessTokenExpiryPeriod       int64  `yaml:"web_access_token_expiry_period"`
	WebRefreshTokenExpiryPeriod      int64  `yaml:"web_refresh_token_expiry_period"`
	SignExpiryPeriod                 int64  `yaml:"sign_expiry_period"`
	AesKey                           string `yaml:"aes_key"`
}

type AccessCountryInfo struct {
	LocationFilePath    string   `yaml:"location_filepath"`
	DisallowedCountries []string `yaml:"disallowed_country"`
	WhiteList           []string `yaml:"white_list"`
	WhiteListMap        map[string]bool
}

type TokenManagerServer struct {
	InternalUri     string `yaml:"api_internal_domain"`
	ExternalUri     string `yaml:"api_external_domain"`
	InternalVersion string `yaml:"internal_ver"`
	ExternalVersion string `yaml:"external_ver"`
}

type PointManagerServer struct {
	InternalUri     string `yaml:"api_internal_domain"`
	ExternalUri     string `yaml:"api_external_domain"`
	InternalVersion string `yaml:"internal_ver"`
	ExternalVersion string `yaml:"external_ver"`
}
type ExternalServer struct {
	InternalUri     string `yaml:"api_internal_domain"`
	ExternalUri     string `yaml:"api_external_domain"`
	InternalVersion string `yaml:"internal_ver"`
	ExternalVersion string `yaml:"external_ver"`
	SecretKey       string `yaml:"secret_key"`
	Network         string `yaml:"network"`
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
	MssqlDBWallet      baseconf.DBAuth `yaml:"mssql_db_wallet"`
	MssqlDBLog         baseconf.DBAuth `yaml:"mssql_db_log"`

	Auth          ApiAuth           `yaml:"api_auth"`
	AccessCountry AccessCountryInfo `yaml:"access_country"`
	TokenManager  ExternalServer    `yaml:"token_manager"`
	PointManager  ExternalServer    `yaml:"point_manager"`
	InnoLog       ExternalServer    `yaml:"inno-log"`
	SuiProv       ExternalServer    `yaml:"sui_prover"`
	SuiEnoki      ExternalServer    `yaml:"sui_enoki"`

	Secret SecretInfo `yaml:"secret"`
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
			currentConfig.AccessCountry.WhiteListMap = make(map[string]bool)
			for _, ip := range currentConfig.AccessCountry.WhiteList {
				slice := strings.Split(ip, ".")
				if len(slice) == 4 && slice[3] == "0" {
					for i := 1; i <= 255; i++ {
						ip = fmt.Sprintf("%v.%v.%v.%v", slice[0], slice[1], slice[2], strconv.Itoa(i))
						currentConfig.AccessCountry.WhiteListMap[ip] = true
					}
				}
				currentConfig.AccessCountry.WhiteListMap[ip] = true
			}
		}
	})

	return currentConfig
}
