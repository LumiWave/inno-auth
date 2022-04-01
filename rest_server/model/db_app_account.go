package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Auth_Members    = "[dbo].[USPAU_Auth_Members]"
	USPAU_Verify_Accounts = "[dbo].[USPAU_Verify_Accounts]"
)

// 앱을 통한 인증 (앱 로그인)
func (o *DB) AuthMembers(account *context.Account, payload *context.Payload) (*context.RespAuthMember, error) {
	resp := new(context.RespAuthMember)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.GetDB().QueryContext(contextR.Background(), USPAU_Auth_Members,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("AppID", payload.AppID),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		sql.Named("MUID", sql.Out{Dest: &resp.MUID}),
		sql.Named("DatabaseID", sql.Out{Dest: &resp.DataBaseID}),
		&returnValue)
	payload.LoginType = context.AppAccountLogin

	defer rows.Close()

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

	if returnValue != 1 {
		return nil, err
	}

	return resp, err
}

func (o *DB) VerfiyAccounts(innoUID string) (bool, error) {
	var returnValue orginMssql.ReturnStatus
	var isExists bool
	_, err := o.MssqlAccountRead.GetDB().QueryContext(contextR.Background(), USPAU_Verify_Accounts,
		sql.Named("InnoUID", innoUID),
		sql.Named("IsExists", sql.Out{Dest: &isExists}),
		&returnValue)
	if err != nil {
		log.Errorf("USPAU_Verify_Accounts QueryContext: %v", err)
		return false, err
	}

	if returnValue != 1 {
		log.Errorf("USPAU_Verify_Accounts returnvalue: %v", returnValue)
		return false, err
	}

	return isExists, err
}
