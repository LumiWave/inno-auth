package commonapi

import (
	"bytes"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/baseapi"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

// [INT] 새 지갑 주소 생성 요청
func GetTokenAddressNew(reqNewWallet *context.ReqNewWallet) (*context.RespNewWallet, error) {
	apiInfo := context.ApiList[context.Api_token_address_new]
	apiInfo.Uri = fmt.Sprintf(apiInfo.Uri, config.GetInstance().TokenManager.Uri)

	apiResp, err := baseapi.HttpCall(apiInfo.Uri, "", "GET", bytes.NewBuffer(nil), reqNewWallet)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if apiResp.Return != 0 {
		// token-manager api error

		return nil, err
	}

	respValue := apiResp.Value.(map[string]interface{})
	resp := new(context.RespNewWallet)
	resp.Symbol = respValue["symbol"].(string)
	resp.Address = respValue["address"].(string)
	return resp, nil
}
