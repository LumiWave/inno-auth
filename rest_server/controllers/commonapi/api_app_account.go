package commonapi

import (
	"net/http"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/inno_log"
	"github.com/ONBUFF-IP-TOKEN/baseInnoClient/point_manager"
	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/commonapi/inner"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostAppAccountLogin(c echo.Context, params *context.Account) error {
	resp := new(base.BaseResponse)
	resp.Success()
	respAccountLogin := new(context.RespAccountLogin)

	// 1. InnoUID 검증
	if isExists, isBlocked, err := model.GetDB().VerfiyAccounts(params.InnoUID); !isExists || isBlocked || err != nil {
		if !isExists {
			log.Errorf("VerifyAccounts err(%v), inno_uid(%v), isExists(%v)", err, params.InnoUID, isExists)
			resp.SetReturn(resultcode.Result_Auth_Invalid_InnoUID)
		} else if isBlocked {
			log.Errorf("VerifyAccounts err(%v), inno_uid(%v), isBlocked(%v)", err, params.InnoUID, isBlocked)
			resp.SetReturn(resultcode.Result_Auth_Blocked_InnoUID)
		} else {
			log.Errorf("VerifyAccounts err(%v), inno_uid(%v), isExists(%v), isBlocked(%v)", err, params.InnoUID, isExists, isBlocked)
			resp.SetReturn(resultcode.Result_DBError)
		}
		return c.JSON(http.StatusOK, resp)
	}

	// 2. 인증 프로시저 호출 (신규 유저, 기존 유저를 체크)
	ctx := base.GetContext(c).(*context.InnoAuthContext)
	if respAuthMember, err := model.GetDB().AuthMembers(params, ctx.Payload); err != nil || respAuthMember == nil {
		log.Errorf("%v", err)
		resp.SetReturn(resultcode.Result_Procedure_Auth_Members)
		return c.JSON(http.StatusOK, resp)
	} else {
		// 3. 신규/기존 유저에 따른 분기 처리
		if respAuthMember.IsJoined {
			// 신규 유저
			// 3-1. [point-manager] 멤버 등록
			if pointList, err := inner.PointMemberRegister(respAuthMember.AUID, respAuthMember.MUID, ctx.Payload.AppID, respAuthMember.DataBaseID); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
				return c.JSON(http.StatusOK, resp)
			} else {
				if len(pointList) == 0 {
					log.Errorf("PointList len=0 err : %v", err)
					resp.SetReturn(resultcode.Result_Api_PointList_Empty)
					return c.JSON(http.StatusOK, resp)
				}
				respAccountLogin.PointList = pointList
			}
		} else {
			// 기존 유저
			// 3-1. [point-manager] 포인트 수량 정보 요청
			if pointList, err := inner.GetPointApp(ctx.Payload.AppID, respAuthMember.MUID, respAuthMember.DataBaseID); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Get_Point_App)
				return c.JSON(http.StatusOK, resp)
			} else {
				// 3-2. [point-manager] pointList가 비어있으면(가입이 안되어있으면) 포인트 멤버 다시 가입
				if len(pointList) == 0 {
					if pointList, err := inner.PointMemberRegister(respAuthMember.AUID, respAuthMember.MUID, ctx.Payload.AppID, respAuthMember.DataBaseID); err != nil {
						log.Errorf("%v", err)
						resp.SetReturn(resultcode.Result_Api_Post_Point_Member_Register)
						return c.JSON(http.StatusOK, resp)
					} else {
						for _, point := range pointList {
							respAccountLogin.PointList = append(respAccountLogin.PointList, point_manager.Point{
								PointID:  point.PointID,
								Quantity: point.Quantity,
							})
						}
					}
				} else {
					for _, point := range pointList {
						respAccountLogin.PointList = append(respAccountLogin.PointList, point_manager.Point{
							PointID:  point.PointID,
							Quantity: point.Quantity,
						})
					}
				}
			}
		}

		// 4. Base Coin의 지갑이 없으면 생성
		if len(respAuthMember.BaseCoinList) > 0 {
			// 4-1. [token-manager] 지갑 생성
			walletInfo, err := inner.TokenAddressNew(respAuthMember.BaseCoinList, params.InnoUID)
			if err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Api_Get_Token_Address_New)
				return c.JSON(http.StatusOK, resp)
			}

			// 4-2. [DB] 지갑 생성 프로시저 호출
			if err := model.GetDB().AddAccountBaseCoins(respAuthMember.AUID, walletInfo); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Base_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}
		}

		// 5. Auth-Members 프로시저에서 내 App에서 사용할 CoinList가 존재하면 지갑 생성
		if len(respAuthMember.AppCoinIDList) > 0 {
			// 5-1. [DB] 사용자 코인 등록 프로시저 호출
			if err := model.GetDB().AddAccountCoins(respAuthMember.AUID, respAuthMember.AppCoinIDList); err != nil {
				log.Errorf("%v", err)
				resp.SetReturn(resultcode.Result_Procedure_Add_Account_Coins)
				return c.JSON(http.StatusOK, resp)
			}
		}

		// 6. [DB] 사용자 로그 등록
		inner.PostMemberAuthLog(&inno_log.MemberAuthLog{
			LogDt:      time.Now().Format("2006-01-02 15:04:05.000"),
			LogID:      context.MemberAuthLog_Auth,
			AUID:       respAuthMember.AUID,
			InnoUID:    params.InnoUID,
			MUID:       respAuthMember.MUID,
			AppID:      ctx.Payload.AppID,
			DataBaseID: respAuthMember.DataBaseID,
		}, respAuthMember.IsJoined)

		// 7. 같이 데이터를 담아서 게임서버로 전달해줌.
		respAccountLogin.MemberInfo.AUID = respAuthMember.AUID
		respAccountLogin.MemberInfo.IsJoined = respAuthMember.IsJoined
		respAccountLogin.MemberInfo.MUID = respAuthMember.MUID
		respAccountLogin.MemberInfo.DataBaseID = respAuthMember.DataBaseID
		resp.Value = respAccountLogin
	}

	return c.JSON(http.StatusOK, resp)
}
