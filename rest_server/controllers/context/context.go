package context

import (
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
)

type LoginType int

const (
	NoneLogin LoginType = iota
	CpLogin
	AppLogin
	AppAccountLogin
	WebAccountLogin
)

var LoginTypeText = map[LoginType]string{
	CpLogin:         "CP",
	AppLogin:        "APP",
	AppAccountLogin: "APPACCOUNT",
	WebAccountLogin: "WEBACCOUNT",
}

type Payload struct {
	CompanyID int64     `json:"company_id,omitempty"`
	AppID     int64     `json:"app_id,omitempty"`
	LoginType LoginType `json:"login_type,omitempty"`
	Uuid      string    `json:"uuid,omitempty"`
	IsEnabled bool      `json:"is_enabled,omitempty"`
	InnoUID   string    `json:"inno_uid,omitempty"`
	AUID      int64     `json:"au_id,omitempty"`
}

// InnoAuthServerContext API의 Request Context
type InnoAuthContext struct {
	*base.BaseContext
	Payload *Payload
}

// NewInnoAuthServerContext 새로운 InnoAuthServer Context 생성
func NewInnoAuthServerContext(baseCtx *base.BaseContext) interface{} {
	if baseCtx == nil {
		return nil
	}

	ctx := new(InnoAuthContext)
	ctx.BaseContext = baseCtx

	return ctx
}

type ReqGetInnoUID struct {
	InnoUID string `json:"inno_uid" query:"inno_uid"`
}

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
}

func (o *InnoAuthContext) SetAuthContext(payload *Payload) {
	o.Payload = payload
}

func MakeDt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}

func GetTokenExpiryperiod(loginType LoginType) (int64, int64) {
	confAuth := config.GetInstance().Auth
	switch loginType {
	case AppLogin:
		return confAuth.AppAccessTokenExpiryPeriod * int64(time.Hour), confAuth.AppRefreshTokenExpiryPeriod * int64(time.Hour)
	case WebAccountLogin:
		return confAuth.WebAccessTokenExpiryPeriod * int64(time.Hour), confAuth.WebRefreshTokenExpiryPeriod * int64(time.Hour)
	}
	return 0, 0
}
