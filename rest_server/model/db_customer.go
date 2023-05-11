package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPWA_Get_CustomerAccounts_By_AccessID = "[dbo].[USPWA_Get_CustomerAccounts_By_AccessID]"
)

// 인증 서버 접근 (앱 로그인/가입)
func (o *DB) GetCustomerAccountsByAccountID(access *context.CustomerAccess) (*context.CustomerPayload, error) {
	payload := new(context.CustomerPayload)

	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlWallet.QueryContext(contextR.Background(), USPWA_Get_CustomerAccounts_By_AccessID,
		sql.Named("AccessID", access.AccessID),
		sql.Named("AccessPW", access.AccessPW),
		sql.Named("AccountID", sql.Out{Dest: &payload.AccountID}),
		sql.Named("CustomerID", sql.Out{Dest: &payload.CustomerID}),
		&returnValue)
	payload.LoginType = context.CustomerLogin

	if rows != nil {
		defer rows.Close()
	}

	if returnValue != 1 {
		return nil, err
	}

	return payload, err
}
