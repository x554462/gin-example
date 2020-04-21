package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
)

type AdminUserDao struct {
	dao.Dao
}

func NewAdminUserDao(ds *dao.DaoSession) *AdminUserDao {
	return ds.GetDao(&model.AdminUser{}, &AdminUserDao{}).(*AdminUserDao)
}

func (this *AdminUserDao) GetRole(adminUser *model.AdminUser) *model.AdminRole {
	adminRoleD := NewAdminRoleDao(this.GetDaoSession())
	return adminRoleD.Select(false, adminUser.RoleId).(*model.AdminRole)
}
