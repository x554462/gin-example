package dao

import (
	"github.com/didi/gendry/builder"
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
	"github.com/x554462/gin-example/middleware/mango/library/database/sqldb"
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

func (this *AdminUserDao) SelectByPage(pageNum, pageSize uint) []dao.ModelInterface {
	cond, vals, err := builder.BuildSelect(this.GetTableName(), map[string]interface{}{
		"_orderby": "id asc",
		"_limit":   []uint{(pageNum - 1) * pageSize, pageSize},
	}, nil)
	this.CheckError(err)
	rows, err := sqldb.GetSlaveDB().Query(cond, vals...)
	this.CheckError(err)
	defer rows.Close()
	return this.CreateMulti(rows)
}
