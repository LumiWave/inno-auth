package model

import (
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Add_AccountCoins     = "[dbo].[USPAU_Add_AccountCoins]"
	USPAU_Add_AccountBaseCoins = "[dbo].[USPAU_Add_AccountBaseCoins]"
)

const (
	TVP_AccountCoins     = "dbo.TVP_AccountCoins"
	TVP_AccountBaseCoins = "dbo.TVP_AccountBaseCoins"
)

// 사용자 코인 등록 [신규가입 시 호출]
func (o *DB) AddAccountCoins(auid int64, coinIDList []int64) error {
	execTvp := "exec " + USPAU_Add_AccountCoins + " @AUID, @TVP;"

	var tableData []context.AccountCoin

	for _, id := range coinIDList {
		data := &context.AccountCoin{
			CoinID:   id,
			Quantity: 0,
		}
		tableData = append(tableData, *data)
	}

	tvpType := orginMssql.TVP{
		TypeName: TVP_AccountCoins,
		Value:    tableData,
	}
	_, err := o.MssqlAccountAll.Exec(execTvp,
		sql.Named("AUID", auid),
		sql.Named("TVP", tvpType))
	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}

// 지갑 생성 [신규가입 시 호출]
func (o *DB) AddAccountBaseCoins(auid int64, walletInfo []*context.WalletInfo) error {
	execTvp := "exec " + USPAU_Add_AccountBaseCoins + " @AUID, @TVP;"

	var tableData []context.AccountBaseCoin

	for _, wallet := range walletInfo {
		data := &context.AccountBaseCoin{
			BaseCoinID:    wallet.BaseCoinID,
			WalletAddress: wallet.Address,
			PrivateKey:    wallet.PrivateKey,
		}
		tableData = append(tableData, *data)
	}

	tvpType := orginMssql.TVP{
		TypeName: TVP_AccountBaseCoins,
		Value:    tableData,
	}
	_, err := o.MssqlAccountAll.Exec(execTvp,
		sql.Named("AUID", auid),
		sql.Named("TVP", tvpType))

	if err != nil {
		log.Errorf("%v", err)
		return err
	}

	return nil
}
