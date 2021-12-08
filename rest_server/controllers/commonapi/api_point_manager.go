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
func PostPointMemberRegister(req *context.ReqPointMemberRegister) (*context.RespPointMemberRegister, error) {
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

	respValue := apiResp.Value.(map[string]interface{})
	points := respValue["points"].([]interface{})
	respPointMemberRegister := new(context.RespPointMemberRegister)

	for _, point := range points {
		data := point.(map[string]interface{})
		p := &context.Point{
			PointID:  int(data["point_id"].(float64)),
			Quantity: int(data["quantity"].(float64)),
		}
		respPointMemberRegister.Points = append(respPointMemberRegister.Points, *p)
	}

	return respPointMemberRegister, nil
}
