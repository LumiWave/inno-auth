package context

import (
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
)

type LoginType int

const (
	Cp LoginType = iota
	App
	Member
)

var LoginTypeText = map[LoginType]string{
	Cp:     "CP",
	App:    "APP",
	Member: "MEMBER",
}

// InnoAuthServerContext API의 Request Context
type InnoAuthContext struct {
	*base.BaseContext
	appInfo   *AppInfo
	loginType LoginType
	Uuid      string
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

func (o *InnoAuthContext) SetAppInfo(appInfo *AppInfo, loginType LoginType, uuid string) {
	o.appInfo = appInfo
	o.loginType = loginType
	o.Uuid = uuid
}

func (o *InnoAuthContext) AppInfo() *AppInfo {
	return o.appInfo
}

func MakeDt(data *int64) {
	*data = datetime.GetTS2MilliSec()
}
