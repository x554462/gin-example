package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/x554462/gin-example/app/dao"
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/app/service"
	"github.com/x554462/gin-example/middleware/mango"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

func PostUserLogin(c *gin.Context) {
	ctrl := mango.Default(c)
	var param struct {
		Account  string `validate:"varchar=用户名,5,32"`
		Passport string `validate:"varchar=密码,40,40"`
	}
	ctrl.ParsePost(&param)
	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())
	adminUser := adminUserD.SelectOne(false, map[string]interface{}{
		"account": param.Account,
	}).(*model.AdminUser)
	if !adminUser.VerifyPassport(param.Passport) {
		exception.ThrowMsg("密码错误，请重新输入", exception.ResponseError)
	}
	ctrl.GetSession().Set(service.AdminUserLoginKey, adminUser.Id)
}

func GetUserList(c *gin.Context) {
	ctrl := mango.Default(c)
	var page struct {
		PageName int `validate:"interge"`
		PageSize int
	}
}
