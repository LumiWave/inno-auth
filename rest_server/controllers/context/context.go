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
	AppAccountLogin
	WebAccountLogin
	CustomerLogin
)

var LoginTypeText = map[LoginType]string{
	CpLogin:         "CP",
	AppLogin:        "APP",
	AppAccountLogin: "APPACCOUNT",
	WebAccountLogin: "WEBACCOUNT",
	CustomerLogin:   "CUSTOMER",
}

type Payload struct {
	CompanyID  int64     `json:"company_id,omitempty"`
	AppID      int64     `json:"app_id,omitempty"`
	LoginType  LoginType `json:"login_type,omitempty"`
	Uuid       string    `json:"uuid,omitempty"`
	IsEnabled  bool      `json:"is_enabled,omitempty"`
	InnoUID    string    `json:"inno_uid,omitempty"`
	AUID       int64     `json:"au_id,omitempty"`
	SocialType int64     `json:"social_type,omitempty"`
	IDToken    string    `json:"id_token,omitempty"`

	Salt string `json:"salt,omitempty"` // 내부 생성용
}

type CustomerPayload struct {
	AccountID  int64     `json:"account_id,omitempty"`
	CustomerID int64     `json:"customer_id,omitempty"`
	LoginType  LoginType `json:"login_type,omitempty"`
	Uuid       string    `json:"uuid,omitempty"`
}

// InnoAuthServerContext API의 Request Context
type InnoAuthContext struct {
	*base.BaseContext
	Payload         *Payload
	CustomerPayload *CustomerPayload
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

func (o *InnoAuthContext) SetAuthPayloadContext(payload *Payload) {
	o.Payload = payload
}

func (o *InnoAuthContext) SetAuthCustomerPayloadContext(payload *CustomerPayload) {
	o.CustomerPayload = payload
}

func MakeDt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}

func (o *InnoAuthContext) GetValue() *Payload {
	return o.Payload
}
