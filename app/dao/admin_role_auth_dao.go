package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
)

type AdminRoleAuthDao struct {
	dao.Dao
}

func NewAdminRoleAuthDao(ds *dao.DaoSession) *AdminRoleAuthDao {
	return ds.GetDao(&model.AdminRoleAuth{}, &AdminRoleAuthDao{}).(*AdminRoleAuthDao)
}

func (this *AdminRoleAuthDao) GetRole(adminRoleAuth *model.AdminRoleAuth) *model.AdminRole {
	adminRoleD := NewAdminRoleDao(this.GetDaoSession())
	return adminRoleD.Select(false, adminRoleAuth.RoleId).(*model.AdminRole)
}

func (this *AdminRoleAuthDao) GetAuth(adminRoleAuth *model.AdminRoleAuth) *model.AdminAuth {
	adminAuthD := NewAdminAuthDao(this.GetDaoSession())
	return adminAuthD.Select(false, adminRoleAuth.AuthId).(*model.AdminAuth)
}
