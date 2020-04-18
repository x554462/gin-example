package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/x554462/gin-example/app/api"
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
		r.GET("/api/:id", api.TestGet)
		r.POST("/api/:id", api.TestPost)
	}
	return r
}
