package context

type api_kind int

const (
	Api_get_token_address_new = iota
	Api_post_point_member_register
	Api_get_point_app
)

type ApiInfo struct {
	Uri string
}

var ApiList = map[api_kind]ApiInfo{
	Api_get_token_address_new:      {Uri: "%s/m1.0/token/address/new"},     // token-manager 새 지갑 주소 생성
	Api_post_point_member_register: {Uri: "%s/m1.0/point/member/register"}, // point-manager 멤버 등록
	Api_get_point_app:              {Uri: "%s/m1.0/point/app"},             // point-manager 포인트 조회
}
