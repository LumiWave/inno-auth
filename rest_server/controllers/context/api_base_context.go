package context

type api_kind int

const (
	Api_create_NewWallet = 0
)

type ApiInfo struct {
	Uri string
}

var ApiList = map[api_kind]ApiInfo{
	Api_create_NewWallet: {Uri: "%s/m1.0/token/address/new"},
}
