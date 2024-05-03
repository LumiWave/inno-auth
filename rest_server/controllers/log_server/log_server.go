package log_server

import (
	"github.com/LumiWave/baseInnoClient/context"
	"github.com/LumiWave/baseInnoClient/inno_log"
	"github.com/LumiWave/inno-auth/rest_server/config"
)

var gLogServer *inno_log.Server

func GetInstance() *inno_log.Server {
	return gLogServer
}

func InitInnoLog(conf *config.ServerConfig) error {
	logServerInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			IntHostUri: conf.InnoLog.InternalUri,
			ExtHostUri: conf.InnoLog.ExternalUri,
			IntVer:     conf.InnoLog.InternalVersion,
			ExtVer:     conf.InnoLog.ExternalVersion,
		},
	}

	gLogServer = inno_log.NewServerInfo(logServerInfo)
	return nil
}
