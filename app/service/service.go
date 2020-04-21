package service

import (
	"github.com/x554462/gin-example/app/dao"
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango"
)

const AdminUserLoginKey = "admin_user"

func GetCurrentAdminUser(ctrl *mango.Controller) *model.AdminUser {
	var (
		userId int
		ok     bool
	)
	if userId, ok = ctrl.GetSession().Get(AdminUserLoginKey).(int); !ok {
		//exception.ThrowMsg("unauthorized", exception.UnauthorizedError)
		userId = 1
	}
	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())
	return adminUserD.Select(false, userId).(*model.AdminUser)
}
