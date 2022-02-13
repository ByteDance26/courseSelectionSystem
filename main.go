package main

import (
	"courseSelectionSystem/DB"
	"courseSelectionSystem/middle"
	"courseSelectionSystem/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	DB.MysqlInit()
	DB.NewRedisHelper()

	middle.InitSimpleSessionPool() //中间件 SimpleSessionPool
	r.Use(middle.HandleSimpleSession("camp-session"))

	router.RegisterRouter(r)
	r.Run(":8000")
}
