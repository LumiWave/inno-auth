package model

import (
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

func (o *DB) InsertAccount(account *context.Account) error {
	return nil
}

func (o *DB) SelectGetAccountInfoByASocialUID(SocialUID string) (*context.Account, error) {
	account := new(context.Account)

	return account, nil
}
