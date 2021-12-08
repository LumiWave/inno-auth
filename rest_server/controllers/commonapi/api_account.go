package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
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
	if respAuthMember, err := model.GetDB().AuthMembers(reqAccountLogin, ctx.Payload); err != nil || respAuthMember == nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_Procedure_Auth_Members)
		return c.JSON(http.StatusOK, resp)
	} else {
		// 2. 신규/기존 유저에 따른 분기 처리
		if respAuthMember.IsJoined == 1 {
			// 신규 유저
			// 1. point-manager 멤버 등록
			if respPointMemberRegister, err := PointMemberRegister(respAuthMember.AUID, respAuthMember.MUID, ctx.Payload.AppID, respAuthMember.DataBaseID); err != nil {
				log.Error(err)
				resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
				return c.JSON(http.StatusOK, resp)
			} else {
				respAccountLogin.Points = respPointMemberRegister.Points
			}

			// 2. token-manager에 새 지갑 주소 생성 요청
			addressList, err := TokenAddressNew(respAuthMember.CoinList, reqAccountLogin.Account.SocialID)
			if err != nil {
				log.Error(err)
				resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
				return c.JSON(http.StatusOK, resp)
			}

			// 3. [DB] 지갑 생성 프로시저 호출
			if err := model.GetDB().AddAccountCoins(respAuthMember, addressList); err != nil {
				log.Error(err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}

		} else {
			// 기존 유저
			// 1. token-manager 호출X -> point-manager에 포인트 수량 정보 요청
		}

		// 3. 같이 데이터를 담아서 게임서버로 전달해줌.
		respAccountLogin.MUID = respAuthMember.MUID
		respAccountLogin.DatabaseID = respAuthMember.DataBaseID
		resp.Value = respAccountLogin
	}

	return c.JSON(http.StatusOK, resp)
}

func PointMemberRegister(AUID int, MUID int, AppID int, DataBaseID int) (*context.RespPointMemberRegister, error) {
	reqPointMemberRegister := &context.ReqPointMemberRegister{
		AUID:       AUID,
		MUID:       MUID,
		AppID:      AppID,
		DataBaseID: DataBaseID,
	}
	return PostPointMemberRegister(reqPointMemberRegister)
}

func TokenAddressNew(coinList []context.CoinInfo, nickName string) ([]context.WalletInfo, error) {
	var addressList []context.WalletInfo

	for _, coin := range coinList {
		reqAddressNew := &context.ReqAddressNew{
			Symbol:   coin.CoinName,
			NickName: nickName,
		}
		if resp, err := GetTokenAddressNew(reqAddressNew); err != nil {
			log.Error(err)
			return nil, err
		} else {
			respAddressNew := &context.WalletInfo{
				CoinID:  coin.CoinID,
				Symbol:  coin.CoinName,
				Address: resp.Address,
			}
			addressList = append(addressList, *respAddressNew)
		}
	}

	return addressList, nil
}
