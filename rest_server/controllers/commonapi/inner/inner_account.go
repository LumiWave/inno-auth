package inner

import (
	"errors"
	"unicode/utf8"

	"github.com/ONBUFF-IP-TOKEN/baseapp/auth/inno"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
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

func DecryptInnoUID(innoUID string) string {
	return inno.AESDecrypt(innoUID, []byte(config.GetInstance().Secret.Key), []byte(config.GetInstance().Secret.Iv))
}

func ValidInnoUID(innoUID string) error {
	// Check InnoUID Length
	if len(innoUID) > 64 {
		return errors.New("invalid inno_uid")
	}
	// Verify InnoUID
	decStr := DecryptInnoUID(innoUID)
	if len(decStr) == 0 || !utf8.ValidString(decStr) {
		return errors.New("invalid inno_uid")
	}
	return nil
}
