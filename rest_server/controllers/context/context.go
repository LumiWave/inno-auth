package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
)

type LoginType int

const (
	NoneLogin LoginType = iota
	CpLogin
	AppLogin
	AccountLogin
)

var LoginTypeText = map[LoginType]string{
	CpLogin:      "CP",
	AppLogin:     "APP",
	AccountLogin: "ACCOUNT",
}

type Payload struct {
	CompanyID int       `json:"company_id"`
	AppID     int       `json:"app_id"`
	LoginType LoginType `json:"login_type"`
	Uuid      string    `json:"uuid"`
}

// InnoAuthServerContext API의 Request Context
type InnoAuthContext struct {
	*base.BaseContext
	Payload *Payload
	JwtInfo *JwtInfo
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

// AppendRequestParameter BaseContext 이미 정의되어 있는 ReqeustParameters 배열에 등록
func AppendRequestParameter() {
}

func (o *InnoAuthContext) SetAuthContext(payload *Payload) {
	o.Payload = payload
}

func MakeDt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}
