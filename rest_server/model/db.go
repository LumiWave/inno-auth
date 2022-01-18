package model

import (
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

type DBMeta struct {
	// 소셜 정보
	Socials    map[int64]*context.SocialInfo
	SocialList []*context.SocialInfo
}

type DB struct {
	Mysql *basedb.Mysql
	Mssql *basedb.Mssql
	Cache *basedb.Cache

	DBMeta
}

var gDB *DB

func SetDB(db *basedb.Mssql, cache *basedb.Cache) {
	gDB = &DB{
		Mssql: db,
		Cache: cache,
	}
	gDB.InitMeta()
}

func GetDB() *DB {
	return gDB
}

func (o *DB) InitMeta() {
	o.Socials = make(map[int64]*context.SocialInfo)

	o.GetSocials()
}
