package model

import (
	"strconv"
	"time"

	baseconf "github.com/LumiWave/baseapp/config"
	"github.com/LumiWave/basedb"
	"github.com/LumiWave/baseutil/log"
	"github.com/LumiWave/inno-auth/rest_server/config"
	"github.com/LumiWave/inno-auth/rest_server/controllers/context"
)

type DBMeta struct {
	// 소셜 정보
	Socials    map[int64]*context.SocialInfo
	SocialList []*context.SocialInfo

	AppCoins  map[int64][]*context.AppCoin // 전체 app에 속한 CoinID 정보 : key AppId
	BaseCoins map[int64]*context.BaseCoinInfo
	CoinsMap  map[int64]*context.CoinInfo // 전체 coin 정보 1 : key CoinId
	Coins     context.CoinList            // 전체 coin 정보 2
}

type DB struct {
	MssqlAccountAll  *basedb.Mssql
	MssqlAccountRead *basedb.Mssql
	MssqlWallet      *basedb.Mssql

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

	if err := ConnectAllDB(conf); err != nil {
		log.Errorf("InitDB Error : %v", err)
		return err
	}

	go func() {
		for {
			timer := time.NewTimer(5 * time.Second)
			<-timer.C
			timer.Stop()

			// DB ping 체크 후 오류 시 재 연결
			if db := CheckPingDB(gDB.MssqlAccountAll, conf.MssqlDBAccountAll); db != nil {
				gDB.MssqlAccountAll = db
			}

			if db := CheckPingDB(gDB.MssqlAccountRead, conf.MssqlDBAccountRead); db != nil {
				gDB.MssqlAccountRead = db
			}

			if db := CheckPingDB(gDB.MssqlWallet, conf.MssqlDBWallet); db != nil {
				gDB.MssqlWallet = db
			}
		}
	}()

	gDB.InitMeta()

	return nil
}

func (o *DB) InitMeta() {
	o.GetSocials()
	o.GetBaseCoins()
	o.GetCoins()
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
	return mssqlDB, nil
}

func ConnectAllDB(conf *config.ServerConfig) error {
	var err error
	gDB.MssqlAccountAll, err = gDB.ConnectDB(&conf.MssqlDBAccountAll)
	if err != nil {
		return err
	}

	gDB.MssqlAccountRead, err = gDB.ConnectDB(&conf.MssqlDBAccountRead)
	if err != nil {
		return err
	}

	gDB.MssqlWallet, err = gDB.ConnectDB(&conf.MssqlDBWallet)
	if err != nil {
		return err
	}
	return nil
}

func CheckPingDB(db *basedb.Mssql, conf baseconf.DBAuth) *basedb.Mssql {
	// 연결이 안되어있거나, DB Connection이 끊어진 경우에는 재연결 시도
	if db == nil || !db.Connection.IsConnect {
		var err error
		newDB, err := gDB.ConnectDB(&conf)
		if err == nil {
			log.Errorf("connect DB OK")
		}
		return newDB
	}

	// 연결이 되어있는 상태면 ping
	if db.Connection.IsConnect {
		if err := db.GetDB().Ping(); err != nil {
			// 재시도 횟수
			db.Connection.RetryCount += 1
			log.Errorf("%v DB Ping err RetryCount(%v)", conf.Database, db.Connection.RetryCount)
			// ping 2회 시도해도 안되면 close
			if db.Connection.RetryCount >= 2 {
				db.Connection.IsConnect = false
				// DB Close
				if err = db.GetDB().Close(); err == nil {
					log.Errorf("DB Closed (RetryCount >=2)")
				}
			}
		}
	}
	return nil
}
