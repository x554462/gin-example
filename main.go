package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/x554462/gin-example/app/router"
	"github.com/x554462/gin-example/middleware/mango/library/conf"
	"log"
	"syscall"
)

func main() {

	gin.SetMode(conf.ServerConf.RunMode)

	endless.DefaultReadTimeOut = conf.ServerConf.ReadTimeout
	endless.DefaultWriteTimeOut = conf.ServerConf.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20

	server := endless.NewServer(conf.ServerConf.Addr, router.InitRouter())
	server.BeforeBegin = func(addr string) {
		log.Printf("start http server listening %s, pid is %d", addr, syscall.Getpid())
	}

	server.ListenAndServe()
}
