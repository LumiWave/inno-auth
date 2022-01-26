package inner

import (
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func TokenAddressNew(coinList []context.CoinInfo, nickName string) ([]context.WalletInfo, error) {
	var addressList []context.WalletInfo

	for _, coin := range coinList {
		reqAddressNew := &context.ReqAddressNew{
			Symbol:   coin.CoinName,
			NickName: nickName,
		}
		if resp, err := GetTokenAddressNew(reqAddressNew); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			respAddressNew := &context.WalletInfo{
				CoinID:  coin.CoinID,
				Symbol:  coin.CoinName,
				Address: resp.Address,
			}
			addressList = append(addressList, *respAddressNew)
		}
	}
	return addressList, nil
}
