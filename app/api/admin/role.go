package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/x554462/gin-example/app/dao"
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango"
	"time"
)

func GetRoleList(c *gin.Context) {
	ctrl := mango.Default(c)

	adminRoleD := dao.NewAdminRoleDao(ctrl.GetDaoSession())

	var data = make([]interface{}, 0)
	for _, v := range adminRoleD.SelectMulti(true, map[string]interface{}{}) {
		role := v.(*model.AdminRole)
		data = append(data, map[string]interface{}{
			"id":          role.Id,
			"name":        role.Name,
			"desc":        role.Describe,
			"create_time": time.Unix(role.CreateTime, 0).Format("2006-01-02 15:04:05"),
		})
	}
	ctrl.JsonResponse(map[string]interface{}{"data": data})
}

func GetRoleById(c *gin.Context) {
	ctrl := mango.Default(c)

	roleId := ctrl.GetPar("roleId").Name("角色id").Min(1).Int()

	adminRoleD := dao.NewAdminRoleDao(ctrl.GetDaoSession())

	role := adminRoleD.Select(false, roleId).(*model.AdminRole)
	ctrl.JsonResponse(map[string]interface{}{
		"id":       role.Id,
		"name":     role.Name,
		"describe": role.Describe,
	})
}

func PostRoleAdd(c *gin.Context) {
	ctrl := mango.Default(c)
	var post struct {
		Name     string `json:"name" validate:"varchar=角色名称,3,20"`
		Describe string `json:"describe" validate:"varchar=描述,1,50"`
	}
	ctrl.ParsePost(&post)

	adminRoleD := dao.NewAdminRoleDao(ctrl.GetDaoSession())

	adminRoleD.Insert(map[string]interface{}{
		"name":        post.Name,
		"describe":    post.Describe,
		"create_time": time.Now().Unix(),
	})
}

func PostRoleSave(c *gin.Context) {
	ctrl := mango.Default(c)
	roleId := ctrl.GetPar("roleId").Name("角色id").Min(1).Int()
	var post struct {
		Name     string `json:"name" validate:"varchar=角色名称,3,20"`
		Describe string `json:"describe" validate:"varchar=描述,1,50"`
	}
	ctrl.ParsePost(&post)

	adminRoleD := dao.NewAdminRoleDao(ctrl.GetDaoSession())

	role := adminRoleD.Select(false, roleId).(*model.AdminRole)
	adminRoleD.Update(role, map[string]interface{}{
		"name":     post.Name,
		"describe": post.Describe,
	})
}
