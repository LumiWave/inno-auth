package sui_prover

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/context"
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/sui_prove"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

var gServer *sui_prove.Server

func GetInstance() *sui_prove.Server {
	return gServer
}

func InitSuiProveManager(conf *config.ServerConfig) error {
	serverInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			IntHostUri: conf.SuiProv.InternalUri,
		},
	}

	gServer = sui_prove.NewServerInfo(serverInfo)
	return nil
}
