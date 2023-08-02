package context

type RespMetaInfo struct {
	Socials []*SocialInfo `json:"social"`
}

// social 리스트 정보 요청
type SocialInfo struct {
	SocialType int64  `json:"social_type" validate:"required"`
	SocialName string `json:"social_name" validate:"required"`
}

// /////// BaseCoinInfo
type BaseCoinInfo struct {
	BaseCoinID         int64  `json:"base_coin_id"`
	BaseCoinName       string `json:"base_coin_name"`
	BaseCoinSymbol     string `json:"base_coin_symbol"`
	IsUsedParentWallet bool   `json:"is_used_parent_wallet"`
}
