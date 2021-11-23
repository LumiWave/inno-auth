package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func (o *DB) InsertCP(params *context.CpInfo) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_cp(cp_name, create_dt) output inserted.idx "+
		"VALUES('%v', %v)",
		params.CpName, params.CreateDt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertCP idx:", lastInsertId)

	return nil
}

func (o *DB) SelectGetCPInfo(param interface{}) (*context.CpInfo, error) {
	var sqlQuery string
	switch param.(type) {
	case int, int64:
		sqlQuery = fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_cp WHERE idx=%v", param)
	case string:
		sqlQuery = fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_cp WHERE cp_name='%v'", param)
	}

	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	cp := new(context.CpInfo)

	for rows.Next() {
		if err := rows.Scan(&cp.Idx, &cp.CpName, &cp.CreateDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return cp, err
}

func (o *DB) DeleteCP(params *context.CpInfo) error {
	sqlQuery := fmt.Sprintf("DELETE FROM onbuff_inno.dbo.auth_cp WHERE cp_name='%v'", params.CpName)
	rows, err := o.Mssql.Query(sqlQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()

	return nil
}
