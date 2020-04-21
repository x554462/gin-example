package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/x554462/gin-example/app/dao"
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/app/service"
	"github.com/x554462/gin-example/middleware/mango"
)

func GetMenu(c *gin.Context) {
	ctrl := mango.Default(c)

	ds := ctrl.GetDaoSession()
	userD := dao.NewAdminUserDao(ds)
	authD := dao.NewAdminAuthDao(ds)

	user := service.GetCurrentAdminUser(ctrl)
	role := userD.GetRole(user)

	var data = make(map[int][]interface{})
	for _, v := range authD.SelectByRoleAndType(role, model.AdminAuthTypeMenu) {
		auth := v.(*model.AdminAuth)
		item := map[string]interface{}{
			"name": auth.Name,
			"path": auth.Path,
		}
		if d, ok := data[auth.Id]; ok {
			item["children"] = d
		}
		data[auth.Pid] = append(data[auth.Pid], item)
	}
	if len(data) > 0 {
		ctrl.JsonResponse(map[string]interface{}{"data": data[0]})
	} else {
		ctrl.JsonResponse(map[string]interface{}{"data": []interface{}{}})
	}
}
