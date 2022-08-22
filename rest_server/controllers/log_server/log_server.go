package log_server

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/context"
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/inno_log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
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
