package token

import "github.com/ONBUFF-IP-TOKEN/baseEthereum/ethcontroller"

func VerifySign(walletAddr, msg, signHex string) bool {
	return ethcontroller.VerifySig(walletAddr, msg, signHex)
}
