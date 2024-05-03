package inner

import (
	"github.com/LumiWave/baseInnoClient/token_manager"
	"github.com/LumiWave/inno-auth/rest_server/controllers/token_server"
)

// token-manager 새 지갑 주소 생성 요청
func GetTokenAddressNew(symbol string, nickName string) (*token_manager.RespAddressNew, error) {
	params := &token_manager.ReqAddressNew{
		BaseSymbol: symbol,
		NickName:   nickName,
	}

	resp, err := token_server.GetInstance().GetTokenAddressNew(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
