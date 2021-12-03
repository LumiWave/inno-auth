package model

import (
	contextR "context"
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"

	orginMssql "github.com/denisenkom/go-mssqldb"
)

func (o *DB) GetApplications(accountInfo *context.AccountInfo) (int, int, int, error) {
	var appID, CompanyID int
	var returnValue orginMssql.ReturnStatus
	_, err := o.Mssql.GetDB().QueryContext(contextR.Background(), "[D-INNO-ACCOUNT01].[dbo].[USPAU_Get_Applications]",
		sql.Named("AccessID", accountInfo.AccessID), sql.Named("AccessPW", accountInfo.AccessPW),
		sql.Named("AppID", sql.Out{Dest: &appID}), sql.Named("CompanyID", sql.Out{Dest: &CompanyID}),
		&returnValue)

	return appID, CompanyID, int(returnValue), err
}

func (o *DB) InsertApp(app *context.Application) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_app(app_name, company_id, access_id, access_pw, create_dt) output inserted.idx "+
		"VALUES('%v', %v, '%v', '%v', %v)",
		app.AppName, app.CompanyID, app.Account.AccessID, app.Account.AccessPW, 0)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertApp idx:", lastInsertId)

	return nil
}

func (o *DB) DeleteApp(app *context.Application) error {
	sqlQuery := fmt.Sprintf("DELETE FROM onbuff_inno.dbo.auth_app WHERE app_name='%v'", app.AppName)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()

	return nil
}

func (o *DB) SelectGetExistsAppAccount(Account context.AccountInfo) (*context.ResponseAppInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_app WHERE login_id='%v' AND login_pwd='%v'", Account.AccessID, Account.AccessPW)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	app := new(context.ResponseAppInfo)

	var loginId, loginPwd, createDt string
	for rows.Next() {
		if err := rows.Scan(&app.AppID, &app.AppName, &app.CompanyID, &loginId, &loginPwd, &createDt); err != nil {
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
		if err := rows.Scan(&app.AppID, &app.AppName, &app.CompanyID, &loginId, &loginPwd, &createDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return app, err
}
