package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
)

type AdminRoleDao struct {
	dao.Dao
}

func NewAdminRoleDao(ds *dao.DaoSession) *AdminRoleDao {
	return ds.GetDao(&model.AdminRole{}, &AdminRoleDao{}).(*AdminRoleDao)
}
