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
	USPR_Add_AccountAuthLogs = "[dbo].[USPR_Add_AccountAuthLogs]"
)

func (o *DB) SetLog(eventid int, auid int64, innoID string, socialID string, socialType int64) error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.MssqlLog.GetDB().QueryContext(contextR.Background(), USPR_Add_AccountAuthLogs,
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
		log.Errorf("USPR_Add_AccountAuthLogs returnvalue error : %v", returnValue)
		return errors.New("USPR_Add_AccountAuthLogs returnvalue error " + strconv.Itoa(int(returnValue)))
	}

	return err
}
