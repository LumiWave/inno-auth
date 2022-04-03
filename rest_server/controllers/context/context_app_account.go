package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
)

// 계정 정보
type Account struct {
	InnoUID string `json:"inno_uid" validate:"required"`
}

// 계정 코인량
type AccountCoin struct {
	CoinID   int64   `json:"coin_id"`
	Quantity float64 `json:"quantity"`
}

// 계정 지갑
type AccountBaseCoin struct {
	BaseCoinID    int64  `json:"base_coin_id"`
	WalletAddress string `json:"wallet_address"`
	PrivateKey    string `json:"private_key"`
}

// 계정 로그인 Response
type RespAccountLogin struct {
	MemberInfo MemberInfo `json:"member_info"`
	PointList  []Point    `json:"points"`
}

// 멤버 정보 (로그인 완료 시 리턴)
type MemberInfo struct {
	AUID       int64 `json:"au_id" url:"au_id"`
	MUID       int64 `json:"mu_id" url:"mu_id"`
	DataBaseID int64 `json:"database_id" url:"database_id"`
	IsJoined   bool  `json:"is_joined"`
}

// 포인트 매니저 PointApp API 호출 Request
type ReqPointApp struct {
	AppID      int64 `json:"app_id" url:"app_id"`
	MUID       int64 `json:"mu_id" url:"mu_id"`
	DataBaseID int64 `json:"database_id" url:"database_id"`
}

// 코인 정보
type CoinInfo struct {
	CoinID     int64  `json:"coin_id,omitempty"`
	CoinSymbol string `json:"coin_symbol,omitempty"`
}

// USPAU_Auth_Members 프로시저 Response
type RespAuthMember struct {
	AUID          int64
	MUID          int64
	DataBaseID    int64
	IsJoined      bool
	BaseCoinList  []CoinInfo
	AppCoinIDList []int64
}

// token-manager 새 지갑 주소 생성 Request
type ReqAddressNew struct {
	Symbol   string `json:"symbol" url:"symbol"`
	NickName string `json:"nickname" url:"nickname"`
}

// token-manager 새 지갑 주소 생성 Response
type WalletInfo struct {
	CoinID     int64  `json:"coin_id"`
	Symbol     string `json:"symbol"`
	Address    string `json:"address"`
	PrivateKey string `json:"pk"`
}

// point-manager 멤버 등록 Request
type ReqPointMemberRegister struct {
	AUID       int64 `json:"au_id"`
	MUID       int64 `json:"mu_id"`
	AppID      int64 `json:"app_id"`
	DataBaseID int64 `json:"database_id"`
}

// 포인트 수량
type Point struct {
	PointID  int64 `json:"point_id"`
	Quantity int64 `json:"quantity"`
}

func NewAccount() *Account {
	return new(Account)
}

func (o *Account) CheckValidate() *base.BaseResponse {
	if len(o.InnoUID) == 0 {
		return base.MakeBaseResponse(resultcode.Result_Auth_EmptyInnoUID)
	}
	return nil
}
