package inner

import (
	"github.com/LumiWave/baseInnoClient/point_manager"
	"github.com/LumiWave/inno-auth/rest_server/controllers/point_server"
)

// point-manager 포인트 멤버 등록
func PointMemberRegister(AUID int64, MUID int64, AppID int64, DataBaseID int64) ([]point_manager.Point, error) {
	params := &point_manager.ReqPointMemberRegister{
		AUID:       AUID,
		MUID:       MUID,
		AppID:      AppID,
		DataBaseID: DataBaseID,
	}
	resp, err := point_server.GetInstance().PostPointMemberRegister(params)
	if err != nil {
		return nil, err
	}

	return resp.Value.Points, nil
}

// point-manager 포인트 조회
func GetPointApp(AppID int64, MUID int64, DatabaseID int64) ([]point_manager.Point, error) {
	params := &point_manager.ReqPointApp{
		AppID:      AppID,
		MUID:       MUID,
		DataBaseID: DatabaseID,
	}
	resp, err := point_server.GetInstance().GetPointApp(params)
	if err != nil {
		return nil, err
	}

	return resp.Value.Points, nil
}
