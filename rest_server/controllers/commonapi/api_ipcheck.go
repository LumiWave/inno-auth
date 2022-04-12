package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/ip"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/labstack/echo"
)

func PostIPAccessAllow(c echo.Context, params *context.ReqIPCheck) error {
	resp := new(base.BaseResponse)
	resp.Success()
	conf := config.GetInstance().AccessCountry

	// check white list
	if access := CheckWhiteList(params.Ip); access {
		respIpCheck := &context.RespIPCheck{
			AllowAccess: access,
		}
		resp.Value = respIpCheck
		return c.JSON(http.StatusOK, resp)
	}

	// check country
	if country, err := ip.GetCountryByIp(params.Ip, conf.LocationFilePath); err != nil {
		resp.SetReturn(resultcode.Result_Auth_Invalid_IPAddress)
	} else {
		respIpCheck := &context.RespIPCheck{
			Country:     country,
			AllowAccess: CheckAllowAccess(country, conf.DisallowedCountries),
		}
		resp.Value = respIpCheck
	}

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
