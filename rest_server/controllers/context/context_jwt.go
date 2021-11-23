package context

type JwtInfo struct {
	Idx int64 `json:"idx,omitempty"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuit"`
	AtExpireDt   int64  `json:"access_token_expire_dt"`
	RtExpireDt   int64  `json:"refresh_token_expire_dt"`
}
