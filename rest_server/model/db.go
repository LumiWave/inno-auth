package model

import (
	"github.com/ONBUFF-IP-TOKEN/basedb"
)

type DB struct {
	Mysql *basedb.Mysql
	Mssql *basedb.Mssql
	Cache *basedb.Cache
}

var gDB *DB

func SetDB(db *basedb.Mssql, cache *basedb.Cache) {
	gDB = &DB{
		Mssql: db,
		Cache: cache,
	}
}

func GetDB() *DB {
	return gDB
}
