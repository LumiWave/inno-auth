package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Auth_Accounts = "[dbo].[USPAU_Auth_Accounts]"
)

// Web 가입/로그인
func (o *DB) AuthAccounts(account *context.ReqAccountWeb) (*context.ResAccountWeb, []*context.NeedWallet, error) {
	resp := new(context.ResAccountWeb)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.QueryContext(contextR.Background(), USPAU_Auth_Accounts,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("SocialID", account.SocialID),
		sql.Named("SocialType", account.SocialType),
		sql.Named("EA", account.EA),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	var needWallets []*context.NeedWallet

	for rows.Next() {
		var baseCoinID int64
		if err := rows.Scan(&baseCoinID); err != nil {
			return nil, nil, err
		} else {
			needWallets = append(needWallets, &context.NeedWallet{
				BaseCoinID: baseCoinID,
			})
		}
	}

	if returnValue != 1 {
		return nil, nil, err
	}

	return resp, needWallets, err
}
