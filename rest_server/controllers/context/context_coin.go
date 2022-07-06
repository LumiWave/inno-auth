package context

// 코인 정보
type CoinInfo struct {
	CoinID     int64  `json:"coin_id,omitempty"`
	CoinSymbol string `json:"coin_symbol,omitempty"`
}

// 필요한 지갑 정보
type NeedWallet struct {
	BaseCoinID     int64
	BaseCoinSymbol string
}
