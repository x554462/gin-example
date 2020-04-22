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
		Account  string `json:"account" validate:"varchar=用户名,5,32"`
		Passport string `json:"passport" validate:"varchar=密码,40,40"`
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
	pageNum := ctrl.DefaultQuery("pageNum", "1").Name("页码").Min(1).Max(2000).UInt()
	pageSize := ctrl.DefaultQuery("pageSize", "10").Name("单页数量").Min(5).Max(50).UInt()

	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())

	var data = make([]interface{}, 0)
	for _, v := range adminUserD.SelectByPage(pageNum, pageSize) {
		user := v.(*model.AdminUser)
		data = append(data, map[string]interface{}{
			"id":      user.Id,
			"account": user.Account,
			"name":    user.Name,
			"role":    adminUserD.GetRole(user).Name,
		})
	}

	ctrl.JsonResponse(map[string]interface{}{"data": data})
}

func GetUserById(c *gin.Context) {
	ctrl := mango.Default(c)
	userId := ctrl.GetPar("userId").Name("用户id").Min(1).Int()

	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())

	user := adminUserD.Select(false, userId).(*model.AdminUser)
	ctrl.JsonResponse(map[string]interface{}{
		"id":      user.Id,
		"account": user.Account,
		"name":    user.Name,
		"role_id": user.RoleId,
	})
}

func PostUserAdd(c *gin.Context) {
	ctrl := mango.Default(c)
	var post struct {
		Account  string `json:"account" validate:"varchar=登录账号,5,20"`
		Passport string `json:"passport" validate:"varchar=密码,40,40"`
		Name     string `json:"name" validate:"varchar=账户名称,3,20"`
		Role     int    `json:"role" validate:"integer=角色id,1"`
	}
	ctrl.ParsePost(&post)

	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())

	adminUserD.Insert(map[string]interface{}{
		"account":  post.Account,
		"passport": post.Passport,
		"name":     post.Name,
		"role_id":  post.Role,
	})
}

func PostUserSave(c *gin.Context) {
	ctrl := mango.Default(c)
	userId := ctrl.GetPar("userId").Name("用户id").Min(1).Int()
	var post struct {
		Account  string `json:"account" validate:"varchar=登录账号,5,20"`
		Passport string `json:"passport" validate:"varchar=密码,40,40"`
		Name     string `json:"name" validate:"varchar=账户名称,3,20"`
		Role     int    `json:"role" validate:"integer=角色id,1"`
	}
	ctrl.ParsePost(&post)

	adminUserD := dao.NewAdminUserDao(ctrl.GetDaoSession())

	user := adminUserD.Select(false, userId).(*model.AdminUser)
	adminUserD.Update(user, map[string]interface{}{
		"account":  post.Account,
		"passport": post.Passport,
		"name":     post.Name,
		"role_id":  post.Role,
	})

}
