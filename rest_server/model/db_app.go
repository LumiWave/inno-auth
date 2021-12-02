package model

import (
	contextR "context"
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"

	orginMssql "github.com/denisenkom/go-mssqldb"
)

func (o *DB) GetApplications(accountInfo *context.AccountInfo) (int64, int64, int64, error) {
	var appID, CompanyID int64
	var returnValue orginMssql.ReturnStatus
	_, err := o.Mssql.GetDB().QueryContext(contextR.Background(), "[D-INNO-ACCOUNT01].[dbo].[USPAU_Get_Applications]",
		sql.Named("AccessID", accountInfo.LoginId), sql.Named("AccessPW", accountInfo.LoginPwd),
		sql.Named("AppID", sql.Out{Dest: &appID}), sql.Named("CompanyID", sql.Out{Dest: &CompanyID}),
		&returnValue)

	return appID, CompanyID, int64(returnValue), err
}

func (o *DB) InsertApp(appInfo *context.AppInfo) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_app(app_name, cp_idx, login_id, login_pwd, create_dt) output inserted.idx "+
		"VALUES('%v', %v, '%v', '%v', %v)",
		appInfo.AppName, appInfo.CpIdx, appInfo.Account.LoginId, appInfo.Account.LoginPwd, appInfo.CreateDt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertApp idx:", lastInsertId)

	return nil
}

func (o *DB) DeleteApp(appInfo *context.AppInfo) error {
	sqlQuery := fmt.Sprintf("DELETE FROM onbuff_inno.dbo.auth_app WHERE app_name='%v'", appInfo.AppName)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()

	return nil
}

func (o *DB) SelectGetExistsAppAccount(Account context.AccountInfo) (*context.ResponseAppInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_app WHERE login_id='%v' AND login_pwd='%v'", Account.LoginId, Account.LoginPwd)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	app := new(context.ResponseAppInfo)

	var loginId, loginPwd, createDt string
	for rows.Next() {
		if err := rows.Scan(&app.Idx, &app.AppName, &app.CpIdx, &loginId, &loginPwd, &createDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return app, err
}

func (o *DB) SelectGetAppInfoByAppName(appName string) (*context.ResponseAppInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_app WHERE app_name='%v'", appName)
	rows, err := o.Mssql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	app := new(context.ResponseAppInfo)

	var loginId, loginPwd, createDt string
	for rows.Next() {
		if err := rows.Scan(&app.Idx, &app.AppName, &app.CpIdx, &loginId, &loginPwd, &createDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return app, err
}
