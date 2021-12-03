package context

type AccountInfo struct {
	AccessID string `json:"access_id" validate:"required"`
	AccessPW string `json:"access_pw" validate:"required"`
}
