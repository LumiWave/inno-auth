package internalapi

import (
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi"
	"github.com/labstack/echo"
)

// 고객사 액세스 토큰 검증
func (o *InternalAPI) GetCustomerTokenVerify(c echo.Context) error {
	return commonapi.GetCustomerTokenVerify(c)
}
