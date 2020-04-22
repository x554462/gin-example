package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/x554462/gin-example/app/api/admin"
	_ "github.com/x554462/gin-example/docs"
	"github.com/x554462/gin-example/middleware/mango"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	if gin.IsDebugging() {
		pprof.Register(r)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.Use(mango.New())
	{
		r.GET("/admin/menu", admin.GetMenu)
		r.POST("/admin/user/login", admin.PostUserLogin)
		r.GET("/admin/user/list", admin.GetUserList)
		r.GET("/admin/user/list/:userId", admin.GetUserById)
		r.POST("/admin/user/add", admin.PostUserAdd)
		r.POST("/admin/user/save/:userId", admin.PostUserSave)
		r.GET("/admin/role/list", admin.GetRoleList)
		r.GET("/admin/role/list/:roleId", admin.GetRoleById)
		r.POST("/admin/role/add", admin.PostRoleAdd)
		r.POST("/admin/role/save/:roleId", admin.PostRoleSave)

	}
	return r
}
