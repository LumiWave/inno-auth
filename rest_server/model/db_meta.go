package model

import (
	contextR "context"

	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	orginMssql "github.com/denisenkom/go-mssqldb"
)

const (
	USPAU_Scan_Socials = "[dbo].[USPAU_Scan_Socials]"
)

func (o *DB) GetSocials() error {
	var returnValue orginMssql.ReturnStatus
	rows, err := o.Mssql.GetDB().QueryContext(contextR.Background(), USPAU_Scan_Socials, &returnValue)

	defer rows.Close()

	for rows.Next() {
		social := &context.SocialInfo{}
		if err := rows.Scan(&social.SocialType, &social.SocialName); err != nil {
			log.Errorf("USPAU_Scan_Socials scan error:%v", err)
		} else {
			o.Socials[social.SocialType] = social
			o.SocialsS = append(o.SocialsS, social)
		}
	}

	if returnValue != 1 {
		return err
	}

	return nil
}
