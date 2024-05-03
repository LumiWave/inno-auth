package model

import (
	contextR "context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Add_AccountCoins     = "[dbo].[USPAU_Add_AccountCoins]"
	USPAU_Add_AccountBaseCoins = "[dbo].[USPAU_Add_AccountBaseCoins]"
	USPAU_GetList_AccountCoins = "[dbo].[USPAU_GetList_AccountCoins]"
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

// 계정 코인 조회
func (o *DB) GetListAccountCoins(auid int64) (map[int64]*context.MeCoin, error) {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_GetList_AccountCoins,
		sql.Named("AUID", auid),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		log.Error("USPAU_GetList_AccountCoins QueryContext err : ", err)
		return nil, err
	}

	meCoinMap := map[int64]*context.MeCoin{}
	for rows.Next() {
		meCoin := &context.MeCoin{}
		if err := rows.Scan(&meCoin.CoinID,
			&meCoin.BaseCoinID,
			&meCoin.WalletAddress,
			&meCoin.Quantity,
			&meCoin.TodayAcqQuantity,
			&meCoin.TodayCnsmQuantity,
			&meCoin.TodayAcqExchangeQuantity,
			&meCoin.TodayCnsmExchangeQuantity,
			&meCoin.ResetDate); err != nil {
			log.Errorf("USPAU_GetList_AccountCoins Scan error %v", err)
			return nil, err
		} else {
			meCoin.CoinSymbol = o.CoinsMap[meCoin.CoinID].CoinSymbol
			meCoinMap[meCoin.CoinID] = meCoin
		}
	}

	if returnValue != 1 {
		log.Errorf("USPAU_GetList_AccountCoins returnvalue error : %v", returnValue)
		return nil, errors.New("USPAU_GetList_AccountCoins returnvalue error " + strconv.Itoa(int(returnValue)))
	}
	return meCoinMap, nil
}
