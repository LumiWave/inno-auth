package context

type ReqIPCheck struct {
	Ip string `json:"ip"`
}

type RespIPCheck struct {
	Country     string `json:"country"`
	AllowAccess bool   `json:"allow_access"`
}
