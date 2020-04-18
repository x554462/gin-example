package sqldb

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/x554462/gin-example/middleware/mango/library/conf"
	"log"
	"sync"
)

var (
	masterDbInstance *sql.DB
	masterDbOnce     sync.Once
	slaveDbInstance  *sql.DB
	slaveDbOnce      sync.Once
)

func GetMasterDB() *sql.DB {
	masterDbOnce.Do(func() {
		var (
			dbConf = conf.MasterDatabaseConf
			err    error
		)
		masterDbInstance, err = manager.New(dbConf.Name, dbConf.User, dbConf.Password, dbConf.Host).Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
		).Port(dbConf.Port).Open(true)
		if err != nil {
			log.Fatalln("dbhelper.DbInstanceMaster,", err)
		}
	})
	return masterDbInstance
}

func GetSlaveDB() *sql.DB {
	slaveDbOnce.Do(func() {
		var (
			dbConf = conf.SlaveDatabaseConf
			err    error
		)
		slaveDbInstance, err = manager.New(dbConf.Name, dbConf.User, dbConf.Password, dbConf.Host).Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
		).Port(dbConf.Port).Open(true)
		if err != nil {
			log.Fatalln("dbhelper.DbInstanceSlave,", err)
		}
	})
	return slaveDbInstance
}
