package dao

import (
	"database/sql"
	"github.com/didi/gendry/builder"
	"github.com/x554462/gin-example/middleware/mango/library/database/sqldb"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"github.com/x554462/gin-example/middleware/mango/library/utils"
	"reflect"
)

const defaultTagName = "db"

type Interface interface {
	initDao(model ModelInterface, ds *DaoSession)
	Select(forUpdate bool, indexes ...interface{}) ModelInterface
	Insert(data map[string]interface{}) ModelInterface
	Update(model ModelInterface, data map[string]interface{}) int64
	Remove(model ModelInterface)
	SelectOne(useSlave bool, where map[string]interface{}) ModelInterface
	SelectMulti(useSlave bool, where map[string]interface{}) []ModelInterface
	CreateOne(row *sql.Rows) ModelInterface
	CreateMulti(rows *sql.Rows) []ModelInterface
}

type Dao struct {
	tableName     string
	indexFields   []string
	notFoundError exception.ErrorWrap
	daoSession    *DaoSession
	model         ModelInterface
}

func (d *Dao) GetDaoSession() *DaoSession {
	return d.daoSession
}

func checkError(err error) {
	if err != nil {
		exception.ThrowMsgWithCallerDepth(err.Error(), exception.ModelRuntimeError, 3)
	}
}

func (d *Dao) newEmptyModel() ModelInterface {
	reflectVal := reflect.ValueOf(d.model)
	t := reflect.Indirect(reflectVal).Type()
	vc := reflect.New(t)
	model, ok := vc.Interface().(ModelInterface)
	if !ok {
		exception.ThrowMsg("dao.newEmptyModel error", exception.ModelRuntimeError)
	}
	return model
}

func (d *Dao) getTableName() string {
	return d.tableName
}

func (d *Dao) buildWhere(indexes ...interface{}) map[string]interface{} {
	if len(d.indexFields) != len(indexes) {
		exception.ThrowMsg("dao.buildWhere index number error", exception.ModelRuntimeError)
	}
	where := make(map[string]interface{})
	for i, v := range d.indexFields {
		where[v] = indexes[i]
	}
	return where
}

func (d *Dao) initDao(model ModelInterface, ds *DaoSession) {
	tableName, indexFields, err := model.InitModelInfo()
	if len(indexFields) == 0 {
		exception.ThrowMsg("dao.initDao: model indexFields empty", exception.ModelRuntimeError)
	}
	d.model = model
	d.tableName = tableName
	d.indexFields = indexFields
	d.notFoundError = err
	d.daoSession = ds
}

func (d *Dao) Select(forUpdate bool, indexes ...interface{}) ModelInterface {
	var (
		daoSession = d.GetDaoSession()
		row        *sql.Rows
	)
	cond, vals, err := builder.BuildSelect(d.getTableName(), d.buildWhere(indexes...), nil)
	checkError(err)
	if forUpdate {
		if daoSession.tx == nil {
			exception.ThrowMsg("Attempt to load for update out of transaction", exception.ModelRuntimeError)
		}
		cond = cond + " FOR UPDATE"
		row, err = daoSession.Query(cond, vals...)
	} else {
		obj := d.query(indexes...)
		if obj != nil {
			return obj
		} else if daoSession.tx != nil {
			row, err = daoSession.Query(cond, vals...)
		} else {
			row, err = sqldb.GetSlaveDB().Query(cond, vals...)
		}
	}
	checkError(err)
	defer row.Close()
	return d.CreateOne(row)
}

func (d *Dao) Insert(data map[string]interface{}) ModelInterface {
	cond, vals, err := builder.BuildInsert(d.getTableName(), []map[string]interface{}{data})
	checkError(err)
	result, err := d.GetDaoSession().Exec(cond, vals...)
	checkError(err)
	if affected, _ := result.RowsAffected(); affected != 1 {
		exception.ThrowMsg("dao.Insert error", exception.ModelRuntimeError)
	}
	if len(d.indexFields) == 1 {
		if id, err := result.LastInsertId(); err == nil {
			data[d.indexFields[0]] = id
		}
	}
	var m = d.newEmptyModel()
	checkError(utils.ScanStruct(data, m, defaultTagName))
	d.save(m)
	return m
}

func (d *Dao) Update(model ModelInterface, data map[string]interface{}) int64 {
	cond, vals, err := builder.BuildUpdate(d.getTableName(), d.buildWhere(model.GetIndexValues()...), data)
	checkError(err)
	result, err := d.GetDaoSession().Exec(cond, vals...)
	checkError(err)
	affected, _ := result.RowsAffected()
	if affected == 1 {
		utils.ScanStruct(data, model, defaultTagName)
		d.save(model)
	}
	return affected
}

func (d *Dao) Remove(model ModelInterface) {
	cond, vals, err := builder.BuildDelete(d.getTableName(), d.buildWhere(model.GetIndexValues()...))
	checkError(err)
	_, err = d.GetDaoSession().Exec(cond, vals...)
	checkError(err)
}

func (d *Dao) SelectOne(useSlave bool, where map[string]interface{}) ModelInterface {
	cond, vals, err := builder.BuildSelect(d.getTableName(), where, nil)
	checkError(err)
	var row *sql.Rows
	if useSlave {
		row, err = sqldb.GetSlaveDB().Query(cond, vals...)
	} else {
		row, err = d.GetDaoSession().Query(cond, vals...)
	}
	checkError(err)
	defer row.Close()
	return d.CreateOne(row)
}

func (d *Dao) SelectMulti(useSlave bool, where map[string]interface{}) []ModelInterface {
	cond, vals, err := builder.BuildSelect(d.getTableName(), where, nil)
	checkError(err)
	var row *sql.Rows
	if useSlave {
		row, err = sqldb.GetSlaveDB().Query(cond, vals...)
	} else {
		row, err = d.GetDaoSession().Query(cond, vals...)
	}
	checkError(err)
	defer row.Close()
	return d.CreateMulti(row)
}

func (d *Dao) CreateOne(row *sql.Rows) ModelInterface {
	columns, err := row.Columns()
	checkError(err)
	length := len(columns)
	values := make([]interface{}, length, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}
	for row.Next() {
		err = row.Scan(values...)
		checkError(err)
		mp := make(map[string]interface{})
		for idx, name := range columns {
			mp[name] = *(values[idx].(*interface{}))
		}
		model := d.newEmptyModel()
		checkError(utils.ScanStruct(mp, model, defaultTagName))
		d.save(model)
		return model
	}
	exception.Throw(d.notFoundError)
	return nil
}

func (d *Dao) CreateMulti(rows *sql.Rows) []ModelInterface {
	modelIs := make([]ModelInterface, 0)
	columns, err := rows.Columns()
	checkError(err)
	length := len(columns)
	values := make([]interface{}, length, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}
	for rows.Next() {
		err = rows.Scan(values...)
		checkError(err)
		mp := make(map[string]interface{})
		for idx, name := range columns {
			mp[name] = *(values[idx].(*interface{}))
		}
		model := d.newEmptyModel()
		checkError(utils.ScanStruct(mp, model, defaultTagName))
		d.save(model)
		modelIs = append(modelIs, model)
	}
	return modelIs
}
