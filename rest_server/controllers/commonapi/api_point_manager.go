package commonapi

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/baseapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

// [INT] 멤버 등록
func PostPointMemberRegister(req *context.ReqPointMemberRegister) ([]context.Point, error) {
	apiInfo := context.ApiList[context.Api_post_point_member_register]
	apiInfo.Uri = fmt.Sprintf(apiInfo.Uri, config.GetInstance().PointManager.Uri)

	pbytes, _ := json.Marshal(req)
	buff := bytes.NewBuffer(pbytes)

	apiResp, err := baseapi.HttpCall(apiInfo.Uri, "", "POST", buff, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if apiResp.Return != 0 {
		log.Error(err)
		return nil, err
	}

	return GetParsePoints(apiResp.Value), nil
}

func GetPointApp(MUID int, DatabaseID int) ([]context.Point, error) {
	apiInfo := context.ApiList[context.Api_get_point_app]
	apiInfo.Uri = fmt.Sprintf(apiInfo.Uri, config.GetInstance().PointManager.Uri)

	memberInfo := &context.MemberInfo{
		MUID:       MUID,
		DataBaseID: DatabaseID,
	}

	apiResp, err := baseapi.HttpCall(apiInfo.Uri, "", "GET", bytes.NewBuffer(nil), memberInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if apiResp.Return != 0 {
		// point-manager api error
		return nil, err
	}

	return GetParsePoints(apiResp.Value), nil
}

func GetParsePoints(value interface{}) []context.Point {
	respValue := value.(map[string]interface{})
	points := respValue["points"].([]interface{})
	var pointList []context.Point
	for _, point := range points {
		data := point.(map[string]interface{})
		p := &context.Point{
			PointID:  int(data["point_id"].(float64)),
			Quantity: int(data["quantity"].(float64)),
		}
		pointList = append(pointList, *p)
	}
	return pointList
}
