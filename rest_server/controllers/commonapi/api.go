package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"

	"github.com/labstack/echo"
)

func GetHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func GetVersion(c echo.Context, maxVersion string) error {
	resp := base.BaseResponse{}

	resp.Value = map[string]interface{}{"version": maxVersion,
		"revision": base.AppVersion}
	resp.Success()

	return c.JSON(http.StatusOK, resp)
}

func PreCheck(c echo.Context) base.PreCheckResponse {
	conf := config.GetInstance()
	if err := base.SetContext(c, &conf.Config, context.NewInnoAuthServerContext); err != nil {
		log.Errorf("%v", err)
		return base.PreCheckResponse{
			IsSucceed: false,
		}
	}

	// auth token 검증
	if conf.Auth.AuthEnable {
		author, ok := c.Request().Header["Authorization"]
		if !ok {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		loginType, atClaims, err := auth.GetIAuth().VerifyAccessToken(author[0][7:])
		var payload *context.Payload
		if err != nil {
			// auth token 오류 리턴
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		} else {
			payload = auth.GetIAuth().ParseClaimsToPayload(loginType, context.AccessT, atClaims)
			base.GetContext(c).(*context.InnoAuthContext).SetAuthPayloadContext(payload)
		}

		if loginType == context.CustomerLogin {
			var customerPayload *context.CustomerPayload
			if err != nil {
				// auth token 오류 리턴
				res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)
				return base.PreCheckResponse{
					IsSucceed: false,
					Response:  res,
				}
			} else {
				customerPayload = auth.GetIAuth().ParseClaimsToCustomerPayload(loginType, context.AccessT, atClaims)
				base.GetContext(c).(*context.InnoAuthContext).SetAuthCustomerPayloadContext(customerPayload)
			}
		}
	}

	return base.PreCheckResponse{
		IsSucceed: true,
	}
}

func GetNodeMetric(c echo.Context) error {
	resp := new(base.BaseResponse)
	resp.Success()

	//node := schedule.GetSystemMonitor().GetMetricInfo()
	//resp.Value = node

	return c.JSON(http.StatusOK, resp)
}

func GetInnoUIDInfo(c echo.Context, innoUID string) error {
	resp := new(base.BaseResponse)
	resp.Success()

	resp.Value = inner.DecryptInnoUID(innoUID)

	return c.JSON(http.StatusOK, resp)
}
