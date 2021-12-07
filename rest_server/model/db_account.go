package model

import (
	contextR "context"
	"database/sql"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

func (o *DB) AccountAuthApplication(reqAuthAccountApp *context.ReqAuthAccountApplication, payload *context.Payload) (*context.RespAuthAccountApplication, error) {
	resp := new(context.RespAuthAccountApplication)
	var returnValue orginMssql.ReturnStatus
	rows, err := o.Mssql.GetDB().QueryContext(contextR.Background(), "[D-INNO-ACCOUNT01].[dbo].[USPAU_Auth_AccountApplications]",
		sql.Named("SocialID", reqAuthAccountApp.Account.SocialID), sql.Named("SocialType", reqAuthAccountApp.Account.SocialType),
		sql.Named("AppID", payload.AppID), sql.Named("CompanyID", payload.CompanyID),
		sql.Named("IsJoined", sql.Out{Dest: &resp.IsJoined}), sql.Named("AUID", sql.Out{Dest: &resp.AUID}), sql.Named("DatabaseID", sql.Out{Dest: &resp.DataBaseID}),
		&returnValue)
	payload.LoginType = context.AccountLogin

	// 신규 유저(IsJoined==1)일 경우 CoinID, CoinName을 추가로 전달 받는다.
	if resp.IsJoined == 1 {
		for rows.Next() {
			if err := rows.Scan(&resp.CoinID, &resp.CoinName); err != nil {
				log.Error(err)
				return nil, err
			}
		}
	}

	return resp, err
}
}

func (o *DB) SelectGetAccountInfoByASocialUID(SocialUID string) (*context.Account, error) {
	account := new(context.Account)

	return account, nil
}
