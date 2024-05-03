package model

import (
	contextR "context"

	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_Socials   = "[dbo].[USPAU_Scan_Socials]"
	USPAU_Scan_BaseCoins = "[dbo].[USPAU_Scan_BaseCoins]"
	USPAU_Scan_Coins     = "[dbo].[USPAU_Scan_Coins]"
)

// 전체 소셜 정보 조회
func (o *DB) GetSocials() error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_Scan_Socials, &returnValue)

	if rows != nil {
		defer rows.Close()
	}

	o.Socials = make(map[int64]*context.SocialInfo)

	for rows.Next() {
		social := &context.SocialInfo{}
		if err := rows.Scan(&social.SocialType, &social.SocialName); err != nil {
			log.Errorf("USPAU_Scan_Socials scan error:%v", err)
		} else {
			o.Socials[social.SocialType] = social
			o.SocialList = append(o.SocialList, social)
		}
	}

	if returnValue != 1 {
		return err
	}

	return nil
}

// Base Coins 정보 조회
func (o *DB) GetBaseCoins() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_Scan_BaseCoins, &rs)
	if err != nil {
		log.Error("QueryContext err : ", err)
		return err
	}

	defer rows.Close()

	o.BaseCoins = make(map[int64]*context.BaseCoinInfo)

	for rows.Next() {
		baseCoin := &context.BaseCoinInfo{}
		if err := rows.Scan(&baseCoin.BaseCoinID, &baseCoin.BaseCoinName, &baseCoin.BaseCoinSymbol, &baseCoin.IsUsedParentWallet); err == nil {
			o.BaseCoins[baseCoin.BaseCoinID] = baseCoin
		}
	}

	return nil
}

// 전체 coin info list
func (o *DB) GetCoins() error {
	var rs orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_Scan_Coins, &rs)
	if err != nil {
		log.Error("USPAU_Scan_Coins QueryContext err : ", err)
		return err
	}

	if rows != nil {
		defer rows.Close()
	}

	o.CoinsMap = make(map[int64]*context.CoinInfo)
	o.Coins.Coins = nil

	for rows.Next() {
		coin := &context.CoinInfo{}
		if err := rows.Scan(&coin.CoinID,
			&coin.BaseCoinID,
			&coin.CoinName,
			&coin.CoinSymbol,
			&coin.ContractAddress,
			&coin.Decimal,
			&coin.ExplorePath,
			&coin.IconUrl,
			&coin.DailyLimitedAcqExchangeQuantity,
			&coin.ExchangeFees,
			&coin.IsRechargeable); err == nil {
			o.Coins.Coins = append(o.Coins.Coins, coin)
			o.CoinsMap[coin.CoinID] = coin
		}
	}

	for _, appCoins := range o.AppCoins {
		for _, appCoin := range appCoins {
			for _, coin := range o.Coins.Coins {
				if appCoin.CoinID == coin.CoinID {
					appCoin.CoinName = coin.CoinName
					appCoin.CoinSymbol = coin.CoinSymbol
					appCoin.ContractAddress = coin.ContractAddress
					appCoin.IconUrl = coin.IconUrl
					appCoin.ExchangeFees = coin.ExchangeFees
					break
				}
			}
		}
	}

	return nil
}
