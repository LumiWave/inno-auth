package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/sui_prove"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/sui_prover"
)

func PostSuiProver(ctx *context.InnoAuthContext, params *context.ReqProve) error {
	resp := new(base.BaseResponse)
	resp.Success()

	req := &sui_prove.ReqProve{
		MaxEpoch:                   params.MaxEpoch,
		JwtRandomness:              params.JwtRandomness,
		Jwt:                        params.Jwt,
		KeyClaimName:               params.KeyClaimName,
		ExtendedEphemeralPublicKey: params.ExtendedEphemeralPublicKey,
		Salt:                       params.Salt,
	}

	if res, err := sui_prover.GetInstance().PostProver(req); err != nil {
		log.Errorf("GetNFTOwned err : %v", err)
		resp.SetReturn(resultcode.ResultInternalServerError)
	} else {
		if len(res.Message) > 0 {
			// error
			resp.SetReturn(resultcode.Result_SUI_Prove)
		}
		resp.Value = res

	}
	return ctx.EchoContext.JSON(http.StatusOK, resp)
}
