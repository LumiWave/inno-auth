package model

import (
	contextR "context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Add_AccountAuthLogs = "[dbo].[USPAU_Add_AccountAuthLogs]"
	USPAU_Add_MemberAuthLogs  = "[dbo].[USPAU_Add_MemberAuthLogs]"
)

func (o *DB) AddAccountAuthLog(eventid int, auid int64, innoID string, socialID string, socialType int64) error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlLog.GetDB().QueryContext(contextR.Background(), USPAU_Add_AccountAuthLogs,
		sql.Named("LogID", 4),
		sql.Named("EventID", eventid),
		sql.Named("AUID", auid),
		sql.Named("InnoUID", innoID),
		sql.Named("SocialID", socialID),
		sql.Named("SocialType", socialType),
		&returnValue)

	if rows != nil {
		defer rows.Close()
	}

	if returnValue != 1 {
		log.Errorf("USPAU_Add_AccountAuthLogs returnvalue error : %v", returnValue)
		return errors.New("USPAU_Add_AccountAuthLogs returnvalue error " + strconv.Itoa(int(returnValue)))
	}

	return err
}

func (o *DB) AddMemberAuthLogs(eventid int64, auid int64, innoUID string, muid int64, appID int64, databaseID int64) error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlLog.GetDB().QueryContext(contextR.Background(), USPAU_Add_MemberAuthLogs,
		sql.Named("LogID", 4),
		sql.Named("EventID", eventid),
		sql.Named("AUID", auid),
		sql.Named("InnoUID", innoUID),
		sql.Named("MUID", muid),
		sql.Named("AppID", appID),
		sql.Named("DatabaseID", databaseID),
		&returnValue)

	if err != nil {
		log.Errorf("USPAU_Add_MemberAuthLogs QueryContext error : %v", err)
		return errors.New("USPAU_Add_MemberAuthLogs QueryContext error")
	}

	if rows != nil {
		defer rows.Close()
	}

	if returnValue != 1 {
		log.Errorf("USPAU_Add_MemberAuthLogs returnvalue error : %v", returnValue)
		return errors.New("USPAU_Add_MemberAuthLogs returnvalue error " + strconv.Itoa(int(returnValue)))
	}

	return err
}
