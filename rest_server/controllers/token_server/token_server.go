package token_server

import (
	"github.com/LumiWave/baseInnoClient/context"
	"github.com/LumiWave/baseInnoClient/token_manager"
	"github.com/LumiWave/inno-auth/rest_server/config"
)

var gTokenServer *token_manager.Server

func GetInstance() *token_manager.Server {
	return gTokenServer
}

func InitTokenManager(conf *config.ServerConfig) error {
	tokenServerInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			IntHostUri: conf.TokenManager.InternalUri,
			ExtHostUri: conf.TokenManager.ExternalUri,
			IntVer:     conf.TokenManager.InternalVersion,
			ExtVer:     conf.TokenManager.ExternalVersion,
		},
	}

	gTokenServer = token_manager.NewServerInfo(tokenServerInfo)
	return nil
}
