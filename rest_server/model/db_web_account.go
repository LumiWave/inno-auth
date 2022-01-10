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

func (o *DB) AuthWebAccounts(account *context.ReqAccountWeb) (*context.ResAccountWeb, error) {
	resp := new(context.ResAccountWeb)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.Mssql.GetDB().QueryContext(contextR.Background(), USPAU_Auth_Accounts,
		sql.Named("InnoUID", account.InnoUID),
		sql.Named("SocialID", account.SocialID),
		sql.Named("SocialType", account.SocialType),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}),
		sql.Named("AUID", sql.Out{Dest: &resp.AUID}),
		&returnValue)

	defer rows.Close()

	if returnValue != 1 {
		return nil, err
	}

	return resp, err
}
