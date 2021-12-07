package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"

	orginMssql "github.com/denisenkom/go-mssqldb"
)

func (o *DB) GetApplications(accountInfo *context.Access) (*context.Payload, int, error) {
	payload := new(context.Payload)
	var returnValue orginMssql.ReturnStatus
	_, err := o.Mssql.GetDB().QueryContext(contextR.Background(), "[D-INNO-ACCOUNT01].[dbo].[USPAU_Get_Applications]",
		sql.Named("AccessID", accountInfo.AccessID), sql.Named("AccessPW", accountInfo.AccessPW),
		sql.Named("AppID", sql.Out{Dest: &payload.AppID}), sql.Named("CompanyID", sql.Out{Dest: &payload.CompanyID}),
		&returnValue)
	payload.LoginType = context.AppLogin

	return payload, int(returnValue), err
}
