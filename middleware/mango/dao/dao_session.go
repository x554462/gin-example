package dao

import (
	"database/sql"
	"fmt"
	"github.com/x554462/gin-example/middleware/mango/library/database/sqldb"
	"github.com/x554462/gin-example/middleware/mango/library/logging"
	"reflect"
	"sync"
)

const daoModelLruCacheSize = 50

type DaoSession struct {
	tx            *sql.Tx
	daoMap        sync.Map
	daoModelCache *DaoLruCache
}

func NewDaoSession() *DaoSession {
	return &DaoSession{daoModelCache: newDaoLru(daoModelLruCacheSize)}
}

func (ds *DaoSession) GetDao(model ModelInterface, daoInterface Interface) Interface {
	name, _, _ := model.InitModelInfo()
	value, ok := ds.daoMap.Load(name)
	if !ok {
		subDaoStructPtrValue := reflect.ValueOf(daoInterface)
		subDaoStructPtrType := reflect.Indirect(subDaoStructPtrValue).Type()
		subDaoStruct := reflect.New(subDaoStructPtrType)
		subDaoStructInterface, _ := subDaoStruct.Interface().(Interface)

		subDaoStructInterface.initDao(model, ds)

		value = subDaoStructInterface
		ds.daoMap.Store(name, value)
	}
	return value.(Interface)
}

func (ds *DaoSession) BeginTransaction() {
	if ds.tx == nil {
		var err error
		if ds.tx, err = sqldb.GetMasterDB().Begin(); err != nil {
			logging.Warn(fmt.Sprintf("daoSession.BeginTransaction: %s\n", err.Error()))
		}
	} else {
		logging.Warn("daoSession.BeginTransaction: can not begin tx again")
	}
}

func (ds *DaoSession) RollbackTransaction() {
	if ds.tx != nil {
		if err := ds.tx.Rollback(); err != nil {
			logging.Warn(fmt.Sprintf("daoSession.RollbackTransaction: %s", err.Error()))
		}
		ds.tx = nil
	}
}

func (ds *DaoSession) SubmitTransaction() {
	if ds.tx != nil {
		if err := ds.tx.Commit(); err != nil {
			logging.Warn(fmt.Sprintf("daoSession.SubmitTransaction: %s", err.Error()))
		}
		ds.tx = nil
	}
}

func (ds *DaoSession) Close() {
	if ds.tx != nil {
		ds.RollbackTransaction()
		ds.tx = nil
	}
}

func (ds *DaoSession) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if ds.tx != nil {
		return ds.tx.Query(query, args...)
	}
	return sqldb.GetMasterDB().Query(query, args...)
}

func (ds *DaoSession) Exec(query string, args ...interface{}) (sql.Result, error) {
	if ds.tx != nil {
		return ds.tx.Exec(query, args...)
	}
	return sqldb.GetMasterDB().Exec(query, args...)
}

func (ds *DaoSession) ClearAllCache() {
	ds.daoModelCache.Clear()
}
