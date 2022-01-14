package context

type RespMetaInfo struct {
	Socials []*SocialInfo `json:"social"`
}

// social 리스트 정보 요청
type SocialInfo struct {
	SocialType int    `json:"social_type" validate:"required"`
	SocialName string `json:"social_name" validate:"required"`
}
