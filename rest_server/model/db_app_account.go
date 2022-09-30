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
	rows, err := o.MssqlAccountAll.QueryContext(contextR.Background(), USPAU_Auth_Members,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("AppID", payload.AppID),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		sql.Named("MUID", sql.Out{Dest: &resp.MUID}),
		sql.Named("DatabaseID", sql.Out{Dest: &resp.DataBaseID}),
		&returnValue)
	payload.LoginType = context.AppAccountLogin

	if rows != nil {
		defer rows.Close()
	}

	// 지갑 생성이 안된 Base Coin List를 전달받는다.
	for rows.Next() {
		var baseCoinID int64
		var baseCoinSymbol string
		if err := rows.Scan(&baseCoinID, &baseCoinSymbol); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			baseCoinInfo := &context.NeedWallet{
				BaseCoinID:     baseCoinID,
				BaseCoinSymbol: baseCoinSymbol,
			}
			resp.BaseCoinList = append(resp.BaseCoinList, baseCoinInfo)
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

func (o *DB) VerfiyAccounts(innoUID string) (bool, bool, error) {
	var returnValue orginMssql.ReturnStatus
	var isExists, isBlocked bool
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_Verify_Accounts,
		sql.Named("InnoUID", innoUID),
		sql.Named("IsExists", sql.Out{Dest: &isExists}),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Errorf("USPAU_Verify_Accounts QueryContext: %v", err)
		return false, false, err
	}

	// 50004 는 제재 유저로 에러 로그를 남기지 않음.
	if returnValue == 50004 {
		isBlocked = true
	} else if returnValue != 1 {
		log.Errorf("USPAU_Verify_Accounts returnvalue: %v", returnValue)
		return false, false, err
	}

	return isExists, isBlocked, err
}
