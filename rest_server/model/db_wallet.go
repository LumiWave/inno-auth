package model

import (
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Add_AccountCoins = "[dbo].[USPAU_Add_AccountCoins]"
)

// 지갑 생성 [신규가입 시 호출]
func (o *DB) AddAccountCoins(auid int64, walletInfo []context.WalletInfo) error {
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
	_, err := o.MssqlAccountAll.GetDB().Exec(execTvp,
		sql.Named("AUID", auid),
		sql.Named("TVP", tvpType))
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}
