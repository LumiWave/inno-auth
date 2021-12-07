package context

type api_kind int

const (
	Api_get_token_address_new = iota
	Api_post_point_member_register
)

type ApiInfo struct {
	Uri string
}

var ApiList = map[api_kind]ApiInfo{
	Api_get_token_address_new:      {Uri: "%s/m1.0/token/address/new"},     // token-manager
	Api_post_point_member_register: {Uri: "%s/m1.0/point/member/register"}, // point-manager
}
