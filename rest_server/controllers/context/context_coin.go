package context

import "time"

// 코인 정보
type CoinInfo struct {
	BaseCoinID                      int64   `json:"base_coin_id"`
	CoinID                          int64   `json:"coin_id,omitempty"`
	CoinName                        string  `json:"coin_name"`
	CoinSymbol                      string  `json:"coin_symbol,omitempty"`
	ContractAddress                 string  `json:"contract_address,omitempty"`
	ExplorePath                     string  `json:"explore_path"`
	IconUrl                         string  `json:"icon_url,omitempty"`
	DailyLimitedAcqExchangeQuantity float64 `json:"daily_limited_acq_exchange_quantity"`
	ExchangeFees                    float64 `json:"exchange_fees"`
}

type CoinList struct {
	Coins []*CoinInfo `json:"coins"`
}

// 필요한 지갑 정보
type NeedWallet struct {
	BaseCoinID     int64
	BaseCoinSymbol string
}

type AppCoin struct {
	AppID int64 `json:"app_id"`
	CoinInfo
}

type MeCoin struct {
	CoinID                    int64     `json:"coin_id" query:"coin_id"`
	BaseCoinID                int64     `json:"base_coin_id" query:"base_coin_id"`
	CoinSymbol                string    `json:"coin_symbol" query:"coin_symbol"`
	WalletAddress             string    `json:"wallet_address" query:"wallet_address"`
	Quantity                  float64   `json:"quantity" query:"quantity"`
	TodayAcqQuantity          float64   `json:"today_acq_quantity" query:"today_acq_quantity"`
	TodayCnsmQuantity         float64   `json:"today_cnsm_quantity" query:"today_cnsm_quantity"`
	TodayAcqExchangeQuantity  float64   `json:"today_acq_exchange_quantity" query:"today_acq_exchange_quantity"`
	TodayCnsmExchangeQuantity float64   `json:"today_cnsm_exchange_quantity" query:"today_cnsm_exchange_quantity"`
	ResetDate                 time.Time `json:"reset_date" query:"reset_date"`
}
