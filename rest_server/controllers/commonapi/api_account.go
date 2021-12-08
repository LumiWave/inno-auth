package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/auth"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAccountLogin(c echo.Context, reqAccountLogin *context.ReqAccountLogin) error {
	resp := new(base.BaseResponse)
	resp.Success()
	respAccountLogin := new(context.RespAccountLogin)

	// 1. 인증 프로시저 호출 (신규 유저, 기존 유저를 체크)
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	if respAuth, err := model.GetDB().AuthMembers(reqAccountLogin, ctx.Payload); err != nil || respAuth == nil {
		// 에러
		log.Error(err)
		resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
		return c.JSON(http.StatusOK, resp)
	} else {
		// 2. 신규/기존 유저에 따른 분기 처리
		if respAuth.IsJoined == 1 {
			// 신규 유저
			// 1. point-manager 멤버 등록
			if err := PointMemberRegister(respAuth.AUID, respAuth.MUID, ctx.Payload.AppID, respAuth.DataBaseID); err != nil {
				resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
				return c.JSON(http.StatusOK, resp)
			}
			// 2. token-manager에 새 지갑 주소 생성 요청
			symbolArray := []string{"ETH", respAuth.CoinName}
			//var addressNewData []*context.RespAddressNew
			for _, symbol := range symbolArray {
				respAddressNew, err := TokenAddressNew(symbol, reqAccountLogin.Account.SocialID)
				if err != nil {
					resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
					return c.JSON(http.StatusOK, resp)
				}
				respAccountLogin.WalletAddress = append(respAccountLogin.WalletAddress, *respAddressNew)
			}

			// 3. [DB] 지갑 생성 프로시저 호출

		} else {
			// 기존 유저
			// 1. token-manager 호출X -> point-manager에 포인트 수량 정보 요청
		}

		// 3. 액세스/리프레시 토큰 생성

		// 4. 같이 데이터를 담아서 액세스토큰과 함께 게임서버로 전달해줌.
	}

	return c.JSON(http.StatusOK, resp)
}

func PointMemberRegister(AUID int, MUID int, AppID int, DataBaseID int) error {
	reqPointMemberRegister := &context.ReqPointMemberRegister{
		AUID:       AUID,
		MUID:       MUID,
		AppID:      AppID,
		DataBaseID: DataBaseID,
	}
	return PostPointMemberRegister(reqPointMemberRegister)
}

func TokenAddressNew(CoinName string, NickName string) (*context.RespAddressNew, error) {
	reqAddressNew := &context.ReqAddressNew{
		Symbol:   CoinName,
		NickName: NickName,
	}
	return GetTokenAddressNew(reqAddressNew)
}
