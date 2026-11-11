package commonapi

import (
	"fmt"
	"net/http"

	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/auth"
	"github.com/LumiWave/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/resultcode"

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
		// "Bearer " 뒤에 토큰이 없으면 잘라낼 때 panic 이 나므로 길이를 먼저 확인한다.
		if len(author) == 0 || len(author[0]) <= 7 {
			res := base.MakeBaseResponse(resultcode.Result_Auth_InvalidJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
		}
		accessToken := author[0][7:]
		loginType, atClaims, err := auth.GetIAuth().VerifyAccessToken(accessToken)
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

		// 서명이 맞아도 서버가 발급해 보관 중인 토큰이 아니면 거절한다.
		// 서명키만으로 만든 위조 토큰, 로그아웃된 토큰, 발급한 적 없는 로그인 종류의 토큰이 여기서 걸린다.
		sessionUuid := payload.Uuid
		if loginType == context.CustomerLogin {
			sessionUuid = fmt.Sprintf("%v", atClaims["access_uuid"])
		}
		if err := auth.GetIAuth().CheckSession(loginType, context.AccessT, sessionUuid, payload.InnoUID, accessToken); err != nil {
			res := base.MakeBaseResponse(resultcode.Result_Auth_ExpiredJwt)

			return base.PreCheckResponse{
				IsSucceed: false,
				Response:  res,
			}
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
