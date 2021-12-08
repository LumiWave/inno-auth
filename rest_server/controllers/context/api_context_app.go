package context

type Application struct {
	AppID     int    `json:"app_id" query:"app_id"`
	AppName   string `json:"app_name"`
	CompanyID int    `json:"company_id"`
	Access    Access `json:"access"`
}

type RequestAppLoginInfo struct {
	Access Access `json:"access" validate:"required"`
}

type ResponseAppInfo struct {
	AppID     int    `json:"app_id"`
	CompanyID int    `json:"company_id"`
	AppName   string `json:"app_name"`
}

func NewApplication() *Application {
	return new(Application)
}
