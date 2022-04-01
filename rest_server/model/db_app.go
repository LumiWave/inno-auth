package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"

	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Get_Applications = "[dbo].[USPAU_Get_Applications]"
)

// 인증 서버 접근 (앱 로그인/가입)
func (o *DB) GetApplications(access *context.Access) (*context.Payload, int, error) {
	payload := new(context.Payload)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlAccountRead.GetDB().QueryContext(contextR.Background(), USPAU_Get_Applications,
		sql.Named("AccessID", access.AccessID), sql.Named("AccessPW", access.AccessPW),
		sql.Named("AppID", sql.Out{Dest: &payload.AppID}),
		sql.Named("CompanyID", sql.Out{Dest: &payload.CompanyID}),
		sql.Named("IsEnabled", sql.Out{Dest: &payload.IsEnabled}),
		&returnValue)
	payload.LoginType = context.AppLogin

	if rows != nil {
		defer rows.Close()
	}

	return payload, int(returnValue), err
}
