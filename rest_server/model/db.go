package model

import (
	"strconv"
	"time"

	baseconf "github.com/ONBUFF-IP-TOKEN/baseapp/config"
	"github.com/ONBUFF-IP-TOKEN/basedb"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
)

type DBMeta struct {
	// 소셜 정보
	Socials    map[int64]*context.SocialInfo
	SocialList []*context.SocialInfo
}

type DB struct {
	MssqlAccountAll  *basedb.Mssql
	MssqlAccountRead *basedb.Mssql

	Cache *basedb.Cache

	DBMeta
}

var gDB *DB

func GetDB() *DB {
	return gDB
}

func InitDB(conf *config.ServerConfig) (err error) {
	cache := basedb.GetCache(&conf.Cache)
	gDB = &DB{
		Cache: cache,
	}

	gDB.MssqlAccountAll, err = gDB.ConnectDB(&conf.MssqlDBAccountAll)
	if err != nil {
		return err
	}

	gDB.MssqlAccountRead, err = gDB.ConnectDB(&conf.MssqlDBAccountRead)
	if err != nil {
		return err
	}

	gDB.InitMeta()

	return nil
}

func (o *DB) InitMeta() {
	o.Socials = make(map[int64]*context.SocialInfo)

	o.GetSocials()
}

func (o *DB) ConnectDB(conf *baseconf.DBAuth) (*basedb.Mssql, error) {
	port, _ := strconv.ParseInt(conf.Port, 10, 32)
	mssqlDB, err := basedb.NewMssql(conf.Database, "", conf.ID, conf.Password, conf.Host, int(port),
		conf.ApplicationIntent, conf.Timeout, conf.ConnectRetryCount, conf.ConnectRetryInterval)
	if err != nil {
		log.Errorf("err: %v, val: %v, %v, %v, %v, %v, %v",
			err, conf.Host, conf.ID, conf.Password, conf.Database, conf.PoolSize, conf.IdleSize)
		return nil, err
	}

	idleSize, _ := strconv.ParseInt(conf.IdleSize, 10, 32)
	mssqlDB.GetDB().SetMaxOpenConns(int(idleSize))
	mssqlDB.GetDB().SetMaxIdleConns(int(idleSize))
	mssqlDB.GetDB().SetConnMaxLifetime(1 * time.Hour)

	return mssqlDB, nil
}
