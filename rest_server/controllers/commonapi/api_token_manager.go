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
func GetTokenAddressNew(reqAddressNew *context.ReqAddressNew) (*context.WalletInfo, error) {
	apiInfo := context.ApiList[context.Api_get_token_address_new]
	apiInfo.Uri = fmt.Sprintf(apiInfo.Uri, config.GetInstance().TokenManager.Uri)

	apiResp, err := baseapi.HttpCall(apiInfo.Uri, "", "GET", bytes.NewBuffer(nil), reqAddressNew)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if apiResp.Return != 0 {
		// token-manager api error

		return nil, err
	}

	respValue := apiResp.Value.(map[string]interface{})
	resp := &context.WalletInfo{
		CoinID:  0,
		Symbol:  respValue["symbol"].(string),
		Address: respValue["address"].(string),
	}
	return resp, nil
}
