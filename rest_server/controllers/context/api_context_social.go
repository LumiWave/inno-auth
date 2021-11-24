package context

type SocialInfo struct {
	SocialUID string `json:"social_uid" validate:"required"`
	SocialID  int64  `json:"social_id" validate:"required"`
}
