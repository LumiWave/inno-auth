package context

type api_kind int

const (
	Api_token_address_new = iota
)

type ApiInfo struct {
	Uri string
}

var ApiList = map[api_kind]ApiInfo{
	Api_token_address_new:     {Uri: "%s/m1.0/token/address/new"},     // token-manager
}
