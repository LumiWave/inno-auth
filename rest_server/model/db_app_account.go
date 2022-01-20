package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Auth_Members     = "[dbo].[USPAU_Auth_Members]"
	USPAU_Add_AccountCoins = "[dbo].[USPAU_Add_AccountCoins]"

	TVP_AccountCoins = "dbo.TVP_AccountCoins"
)

func (o *DB) AuthMembers(account *context.Account, payload *context.Payload) (*context.RespAuthMember, error) {
	resp := new(context.RespAuthMember)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.Mssql.GetDB().QueryContext(contextR.Background(), USPAU_Auth_Members,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("AppID", payload.AppID),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		sql.Named("MUID", sql.Out{Dest: &resp.MUID}),
		sql.Named("DatabaseID", sql.Out{Dest: &resp.DataBaseID}),
		&returnValue)
	payload.LoginType = context.AppAccountLogin

	// 신규 유저(IsJoined==1)일 경우 CoinID, CoinName을 추가로 전달 받는다.
	var coinID int64
	var coinName string
	for rows.Next() {
		if err := rows.Scan(&coinID, &coinName); err != nil {
			log.Errorf("%v", err)
			return nil, err
		} else {
			coinInfo := &context.CoinInfo{
				CoinID:   coinID,
				CoinName: coinName,
			}
			resp.CoinList = append(resp.CoinList, *coinInfo)
		}
	}
	defer rows.Close()

	if returnValue != 1 {
		return nil, err
	}

	return resp, err
}

func (o *DB) AddAccountCoins(respAuthMember *context.RespAuthMember, walletInfo []context.WalletInfo) error {
	execTvp := "exec " + USPAU_Add_AccountCoins + " @AUID, @TVP;"

	var tableData []context.AccountCoin

	for _, wallet := range walletInfo {
		data := &context.AccountCoin{
			CoinID:        wallet.CoinID,
			WalletAddress: wallet.Address,
			Quantity:      0,
		}
		tableData = append(tableData, *data)
	}

	tvpType := orginMssql.TVP{
		TypeName: TVP_AccountCoins,
		Value:    tableData,
	}
	_, err := o.Mssql.GetDB().Exec(execTvp,
		sql.Named("AUID", respAuthMember.AUID),
		sql.Named("TVP", tvpType))
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}
