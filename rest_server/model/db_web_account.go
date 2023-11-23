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
func (o *DB) AuthAccounts(account *context.ReqAccountWeb) (*context.ResAccountWeb, error) {
	resp := new(context.ResAccountWeb)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountAll.QueryContext(contextR.Background(), USPAU_Auth_Accounts,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("SocialID", account.SocialID),
		sql.Named("SocialType", account.SocialType),
		sql.Named("EA", account.EA),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		sql.Named("IsMigrated", sql.Out{Dest: &resp.IsMigrated}),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if returnValue != 1 {
		return nil, err
	}

	return resp, err
}
