package model

import (
	contextR "context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_Socials   = "[dbo].[USPAU_Scan_Socials]"
	USPAU_Scan_BaseCoins = "[dbo].[USPAU_Scan_BaseCoins]"
)

// 전체 소셜 정보 조회
func (o *DB) GetSocials() error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.QueryContext(contextR.Background(), USPAU_Scan_Socials, &returnValue)

	if rows != nil {
		defer rows.Close()
	}

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

	eth := config.GetInstance().EthToken
	matic := config.GetInstance().MaticToken
	o.BaseCoins = make(map[int64]*context.BaseCoinInfo)

	for rows.Next() {
		baseCoin := &context.BaseCoinInfo{}
		if err := rows.Scan(&baseCoin.BaseCoinID, &baseCoin.BaseCoinName, &baseCoin.BaseCoinSymbol, &baseCoin.IsUsedParentWallet); err == nil {
			if baseCoin.BaseCoinSymbol == "ETH" {
				baseCoin.IDList = eth.IDList
			} else if baseCoin.BaseCoinSymbol == "MATIC" {
				baseCoin.IDList = matic.IDList
			}
			o.BaseCoins[baseCoin.BaseCoinID] = baseCoin
		}
	}

	return nil
}
