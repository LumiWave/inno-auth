package auth

var gAuth *IAuth

type LoginVD struct {
	WalletAddr string `json:"wallet_address"`
	Date       int64  `json:"date"`
}
