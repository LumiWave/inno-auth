package api_inno_log

type Common struct {
	Return  int    `json:"return"`
	Message string `json:"message"`
}

type AccountAuthLog struct {
	LogDt      string `json:"log_dt"`
	LogID      int64  `json:"log_id"`
	EventID    int64  `json:"event_id"`
	AUID       int64  `json:"au_id"`
	InnoUID    string `json:"inno_id"`
	SocialID   string `json:"social_id"`
	SocialType int64  `json:"social_type"`
}

type AccountCoinLog struct {
	LogDt         string  `json:"log_dt"`
	LogID         int64   `json:"log_id"`
	EventID       int64   `json:"event_id"`
	AUID          int64   `json:"au_id"`
	CoinID        int64   `json:"coin_id"`
	BaseCoinID    int64   `json:"basecoin_id"`
	WalletAddress string  `json:"wallet_address"`
	PreQuantity   float64 `json:"previous_quantity"`
	AdjQuantity   float64 `json:"adjust_quantity"`
	Quantity      float64 `json:"quantity"`
}

type ExchangeGoodsLog struct {
	LogDt            string  `json:"log_dt"`
	LogID            int64   `json:"log_id"`
	EventID          int64   `json:"event_id"`
	AUID             int64   `json:"au_id"`
	MUID             int64   `json:"mu_id"`
	AppID            int64   `json:"app_id"`
	CoinID           int64   `json:"coin_id"`
	BaseCoinID       int64   `json:"basecoin_id"`
	WalletAddress    string  `json:"wallet_address"`
	PreCoinQuantity  float64 `json:"previous_coin_quantity"`
	AdjCoinQuantity  float64 `json:"adjust_coin_quantity"`
	CoinQuantity     float64 `json:"coin_quantity"`
	PointID          int64   `json:"point_id"`
	PrePointQuantity int64   `json:"previous_point_quantity"`
	AdjPointQuantity int64   `json:"adjust_point_quantity"`
	PointQuantity    int64   `json:"point_quantity"`
}

type MemberAuthLog struct {
	LogDt      string `json:"log_dt"`
	LogID      int64  `json:"log_id"`
	EventID    int64  `json:"event_id"`
	AUID       int64  `json:"au_id"`
	InnoUID    string `json:"inno_id"`
	MUID       int64  `json:"mu_id"`
	AppID      int64  `json:"app_id"`
	DataBaseID int64  `json:"database_id"`
}

type MemberPointsLog struct {
	LogDt       string `json:"log_dt"`
	LogID       int64  `json:"log_id"`
	EventID     int64  `json:"event_id"`
	AUID        int64  `json:"au_id"`
	MUID        int64  `json:"mu_id"`
	AppID       int64  `json:"app_id"`
	PointID     int64  `json:"point_id"`
	PreQuantity int64  `json:"previous_quantity"`
	AdjQuantity int64  `json:"adjust_quantity"`
	Quantity    int64  `json:"quantity"`
}
