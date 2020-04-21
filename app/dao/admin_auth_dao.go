package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
	"github.com/x554462/gin-example/middleware/mango/library/database/sqldb"
)

type AdminAuthDao struct {
	dao.Dao
}

func NewAdminAuthDao(ds *dao.DaoSession) *AdminAuthDao {
	return ds.GetDao(&model.AdminAuth{}, &AdminAuthDao{}).(*AdminAuthDao)
}

func (this *AdminAuthDao) SelectByRoleAndType(role *model.AdminRole, typ int8) []dao.ModelInterface {
	sql := "SELECT `admin_auth`.* FROM `admin_auth` INNER JOIN `admin_role_auth` ON `admin_auth`.`id`= `admin_role_auth`.`auth_id` WHERE `admin_role_auth`.`role_id`= ? AND `admin_auth`.`type`= ? ORDER BY `admin_auth`.`depth` DESC, `admin_auth`.`sort_num` DESC"
	rows, _ := sqldb.GetSlaveDB().Query(sql, role.Id, typ)
	defer rows.Close()
	return this.CreateMulti(rows)
}
