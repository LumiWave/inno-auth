package sui_enoki_server

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/context"
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/sui_enoki"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

var gServer *sui_enoki.Server

func GetInstance() *sui_enoki.Server {
	return gServer
}

func InitSuiEnokiManager(conf *config.ServerConfig) error {
	serverInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			ExtHostUri: conf.SuiEnoki.ExternalUri,
			ExtVer:     conf.SuiEnoki.ExternalVersion,
		},
		AuthInfo: context.AuthInfo{
			ApiKey: conf.SuiEnoki.SecretKey,
		},
	}

	gServer = sui_enoki.NewServerInfo(serverInfo)
	return nil
}
