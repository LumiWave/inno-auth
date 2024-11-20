package commonapi

import (
	"encoding/json"
	"net/http"

	"github.com/LumiWave/baseInnoClient/sui_enoki"
	"github.com/LumiWave/baseapp/base"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	"github.com/LumiWave/inno-auth/rest_server/controllers/resultcode"
	"github.com/LumiWave/inno-auth/rest_server/controllers/sui_enoki_server"
	"github.com/labstack/echo"
)

func PostSuiProver(ctx *context.InnoAuthContext, params *context.ReqProve) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// req := &sui_prove.ReqProve{
	// 	MaxEpoch:                   params.MaxEpoch,
	// 	JwtRandomness:              params.JwtRandomness,
	// 	Jwt:                        params.Jwt,
	// 	KeyClaimName:               params.KeyClaimName,
	// 	ExtendedEphemeralPublicKey: params.ExtendedEphemeralPublicKey,
	// 	Salt:                       params.Salt,
	// }

	// if res, err := sui_prover.GetInstance().PostProver(req); err != nil {
	// 	log.Errorf("PostProver err : %v", err)
	// 	resp.SetReturn(resultcode.ResultInternalServerError)
	// } else {
	// 	if len(res.Message) > 0 {
	// 		// error
	// 		resp.SetReturn(resultcode.Result_SUI_Prove)
	// 	}
	// 	resp.Value = res
	// }

	req := &sui_enoki.ReqzkLoginZKP{
		Network:            config.GetInstance().SuiEnoki.Network,
		Randomness:         params.JwtRandomness,
		MaxEpoch:           params.MaxEpoch,
		EphemeralPublicKey: params.EphemeralPublicKey,
	}
	respEnoki, errEnoki, err := sui_enoki_server.GetInstance().PostZkloginZkp(req, params.Jwt)
	if err != nil {
		log.Errorf("PostZkloginZkp err : %v", err)
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else if errEnoki != nil {
		temp, _ := json.Marshal(errEnoki)
		log.Errorf("PostZkloginZkp errEnoki : %v", string(temp))
		resp.SetReturn(resultcode.Result_Auth_SuiEnoki_ZKP_Error)
	} else {
		resp.Value = respEnoki
	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}

func PostSuiProverNonce(c echo.Context, params *context.ReqProveNonce) error {
	resp := new(base.BaseResponse)
	resp.Success()

	req := &sui_enoki.ReqzkLoginNonce{
		Network:            config.GetInstance().SuiEnoki.Network,
		EphemeralPublicKey: params.EphemeralPublicKey,
	}
	respEnoki, errEnoki, err := sui_enoki_server.GetInstance().PostZkloginNonce(req)
	if err != nil {
		log.Errorf("PostZkloginNonce err : %v", err)
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else if errEnoki != nil {
		temp, _ := json.Marshal(errEnoki)
		log.Errorf("PostZkloginNonce errEnoki : %v", string(temp))
		resp.SetReturn(resultcode.Result_Auth_SuiEnoki_ZKP_Error)
	} else {
		resp.Value = respEnoki
	}
	return c.JSON(http.StatusOK, resp)
}
