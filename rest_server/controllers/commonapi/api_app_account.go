package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAppAccountLogin(c echo.Context, account *context.Account) error {
	resp := new(base.BaseResponse)
	resp.Success()
	respAccountLogin := new(context.RespAccountLogin)

	// // 0. InnoUID 검증
	// if isExists, err := model.GetDB().VerfiyAccounts(account.InnoUID); !isExists || err != nil {
	// 	resp.SetReturn(resultcode.Result_Auth_Invalid_InnoUID)
	// 	return c.JSON(http.StatusOK, resp)
	// }

	// 1. 인증 프로시저 호출 (신규 유저, 기존 유저를 체크)
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	if respAuthMember, err := model.GetDB().AuthMembers(account, ctx.Payload); err != nil || respAuthMember == nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_Procedure_Auth_Members)
		return c.JSON(http.StatusOK, resp)
	} else {
		// 2. 신규/기존 유저에 따른 분기 처리
		if respAuthMember.IsJoined {
			// 신규 유저
			// 2-1. point-manager 멤버 등록
			if pointList, err := PointMemberRegister(respAuthMember.AUID, respAuthMember.MUID, ctx.Payload.AppID, respAuthMember.DataBaseID); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
				return c.JSON(http.StatusOK, resp)
			} else {
				respAccountLogin.PointList = pointList
			}
		} else {
			// 기존 유저
			// 2-1. point-manager에 포인트 수량 정보 요청
			if pointList, err := inner.GetPointApp(ctx.Payload.AppID, respAuthMember.MUID, respAuthMember.DataBaseID); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Get_Point_App)
				return c.JSON(http.StatusOK, resp)
			} else {
				// 2-2. pointList가 비어있으면(가입이 안되어있으면) 포인트 멤버 다시 가입
				if len(pointList) == 0 {
					if pointList, err := PointMemberRegister(respAuthMember.AUID, respAuthMember.MUID, ctx.Payload.AppID, respAuthMember.DataBaseID); err != nil {
						log.Errorf("%v", err)
						resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
						return c.JSON(http.StatusOK, resp)
					} else {
						respAccountLogin.PointList = pointList
					}
				} else {
					respAccountLogin.PointList = pointList
				}
			}
		}

		// 3. Base Coin의 지갑이 없으면 생성
		if len(respAuthMember.BaseCoinList) > 0 {
			// 3-1. [token-manager] 지갑 생성
			walletInfo, err := inner.TokenAddressNew(respAuthMember.BaseCoinList, account.InnoUID)
			if err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
				return c.JSON(http.StatusOK, resp)
			}

			// 3-2. [DB] 지갑 생성 프로시저 호출
			if err := model.GetDB().AddAccountBaseCoins(respAuthMember.AUID, walletInfo); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Base_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}
		}

		// 4. Auth-Members 프로시저에서 내 App에서 사용할 CoinList가 존재하면 지갑 생성
		if len(respAuthMember.AppCoinIDList) > 0 {
			// 4-1. [DB] 사용자 코인 등록 프로시저 호출
			if err := model.GetDB().AddAccountCoins(respAuthMember.AUID, respAuthMember.AppCoinIDList); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}
		}

		// 5. 같이 데이터를 담아서 게임서버로 전달해줌.
		respAccountLogin.MemberInfo.AUID = respAuthMember.AUID
		respAccountLogin.MemberInfo.IsJoined = respAuthMember.IsJoined
		respAccountLogin.MemberInfo.MUID = respAuthMember.MUID
		respAccountLogin.MemberInfo.DataBaseID = respAuthMember.DataBaseID
		resp.Value = respAccountLogin
	}

	return c.JSON(http.StatusOK, resp)
}

func PointMemberRegister(AUID int64, MUID int64, AppID int64, DataBaseID int64) ([]context.Point, error) {
	reqPointMemberRegister := &context.ReqPointMemberRegister{
		AUID:       AUID,
		MUID:       MUID,
		AppID:      AppID,
		DataBaseID: DataBaseID,
	}
	return inner.PostPointMemberRegister(reqPointMemberRegister)
}
