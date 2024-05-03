package sui_prover

import (
	"github.com/LumiWave/baseInnoClient/context"
	"github.com/LumiWave/baseInnoClient/sui_prove"
	"github.com/LumiWave/inno-auth/rest_server/config"
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
