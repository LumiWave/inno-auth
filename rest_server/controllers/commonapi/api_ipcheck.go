package commonapi

import (
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/ip"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostIPAccessAllow(c echo.Context, params *context.ReqIPCheck) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance().AccessCountry

	// check white list
	if access := CheckWhiteList(params.Ip); access {
		respIpCheck := &context.RespIPCheck{
			Country:     "WL",
			AllowAccess: access,
			SwapEnable:  true,
		}
		resp.Value = respIpCheck
		return c.JSON(http.StatusOK, resp)
	}

	// check country
	if country, err := ip.GetCountryByIp(params.Ip, conf.LocationFilePath); err != nil {
		resp.SetReturn(resultcode.Result_Auth_Invalid_IPAddress)
	} else {
		// swap 가능 상태 체크
		bSwapEnable, _ := model.GetDB().GetSwapEnable()
		respIpCheck := &context.RespIPCheck{
			Country:     country,
			AllowAccess: CheckAllowAccess(country, conf.DisallowedCountries),
			SwapEnable:  bSwapEnable,
		}

		resp.Value = respIpCheck
	}

	return c.JSON(http.StatusOK, resp)
}

func PutSwapEnable(c echo.Context, params *context.ReqSwapEnable) error {
	resp := new(base.BaseResponse)
	resp.Success()

	model.GetDB().SetSwapEnable(params.SwapEnable)

	return c.JSON(http.StatusOK, resp)
}

func CheckWhiteList(ip string) bool {
	conf := config.GetInstance().AccessCountry

	if _, ok := conf.WhiteListMap[ip]; ok {
		return true
	}

	return false
}

func CheckAllowAccess(country string, disAllowedCountries []string) bool {
	for _, value := range disAllowedCountries {
		if value == country {
			return false
		}
	}

	return true
}
