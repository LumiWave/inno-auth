package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Auth_Members = "[dbo].[USPAU_Auth_Members]"
)

// 앱을 통한 인증 (앱 로그인)
func (o *DB) AuthMembers(account *context.Account, payload *context.Payload) (*context.RespAuthMember, error) {
	resp := new(context.RespAuthMember)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(contextR.Background(), USPAU_Auth_Members,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("AppID", payload.AppID),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		sql.Named("MUID", sql.Out{Dest: &resp.MUID}),
		sql.Named("DatabaseID", sql.Out{Dest: &resp.DataBaseID}),
		&returnValue)
	payload.LoginType = context.AppAccountLogin

	// 지갑 생성이 안된 Base Coin List를 전달받는다.
	for rows.Next() {
		var baseCoinID int64
		var baseCoinSymbol string
		if err := rows.Scan(&baseCoinID, &baseCoinSymbol); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			baseCoinInfo := &context.CoinInfo{
				CoinID:     baseCoinID,
				CoinSymbol: baseCoinSymbol,
			}
			resp.BaseCoinList = append(resp.BaseCoinList, *baseCoinInfo)
		}
	}
	rows.NextResultSet()

	// 사용자 코인 등록이 안된 Base Coin List를 전달받는다.
	// App CoinList
	for rows.Next() {
		var coinID int64
		if err := rows.Scan(&coinID); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			resp.AppCoinIDList = append(resp.AppCoinIDList, coinID)
		}
	}

	defer rows.Close()

	if returnValue != 1 {
		return nil, err
	}

	return resp, err
}
