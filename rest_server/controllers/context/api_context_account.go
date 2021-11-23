package context

type AccountInfo struct {
	LoginId  string `json:"login_id" query:"login_id" validate:"required"`
	LoginPwd string `json:"login_pwd" query:"login_pwd" validate:"required"`
}
