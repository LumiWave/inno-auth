package model

import (
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func (o *DB) InsertMember(memberInfo *context.MemberInfo) error {
	sqlQuery := fmt.Sprintf("INSERT INTO onbuff_inno.dbo.auth_member(member_id, app_idx, social_uid, social_id, create_dt) output inserted.idx "+
		"VALUES('%v', %v, '%v', %v, %v)",
		memberInfo.MemberID, memberInfo.AppIdx, memberInfo.Social.SocialUID, memberInfo.Social.SocialID, memberInfo.CreateDt)

	var lastInsertId int64
	err := o.Mssql.QueryRow(sqlQuery, &lastInsertId)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug("InsertMember idx:", lastInsertId)

	return nil
}

func (o *DB) SelectGetMemberInfoByASocialUID(SocialUID string) (*context.MemberInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM onbuff_inno.dbo.auth_member WHERE social_uid='%v'", SocialUID)
	rows, err := o.Mssql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	member := new(context.MemberInfo)

	for rows.Next() {
		if err := rows.Scan(&member.Idx, &member.MemberID, &member.AppIdx, &member.Social.SocialUID, &member.Social.SocialID, &member.CreateDt); err != nil {
			log.Error(err)
			return nil, err
		}
	}

	return member, err
}
