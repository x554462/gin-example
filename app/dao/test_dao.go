package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
)

type TestDao struct {
	dao.Dao
}

func NewTestDao(ds *dao.DaoSession) *TestDao {
	dao.NewDaoSession()
	return ds.GetDao(&model.Test{}, &TestDao{}).(*TestDao)
}
