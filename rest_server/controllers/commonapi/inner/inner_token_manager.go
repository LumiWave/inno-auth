package inner

import (
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/token_manager"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/token_server"
)

// token-manager 새 지갑 주소 생성 요청
func GetTokenAddressNew(symbol string, nickName string) (*token_manager.RespAddressNew, error) {
	params := &token_manager.ReqAddressNew{
		Symbol:   symbol,
		NickName: nickName,
	}

	resp, err := token_server.GetInstance().GetTokenAddressNew(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func TokenAddressNew(coinList []*context.NeedWallet, nickName string) ([]*context.WalletInfo, error) {
	var addressList []*context.WalletInfo

	for _, coin := range coinList {
		if resp, err := GetTokenAddressNew(coin.BaseCoinSymbol, nickName); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			respAddressNew := &context.WalletInfo{
				BaseCoinID:     coin.BaseCoinID,
				BaseCoinSymbol: coin.BaseCoinSymbol,
				Address:        resp.Value.Address,
				PrivateKey:     resp.Value.PrivateKey,
			}
			addressList = append(addressList, respAddressNew)
		}
	}
	return addressList, nil
}
