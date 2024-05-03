package point_server

import (
	"github.com/LumiWave/baseInnoClient/context"
	"github.com/LumiWave/baseInnoClient/point_manager"
	"github.com/LumiWave/inno-auth/rest_server/config"
)

var gPointServer *point_manager.Server

func GetInstance() *point_manager.Server {
	return gPointServer
}

func InitPointManager(conf *config.ServerConfig) error {
	pointServerInfo := &context.ServerInfo{
		HostInfo: context.HostInfo{
			IntHostUri: conf.PointManager.InternalUri,
			ExtHostUri: conf.PointManager.ExternalUri,
			IntVer:     conf.PointManager.InternalVersion,
			ExtVer:     conf.PointManager.ExternalVersion,
		},
	}

	gPointServer = point_manager.NewServerInfo(pointServerInfo)
	return nil
}
