package main

import (
	"courseSelectionSystem/DB"
	"courseSelectionSystem/middle"
	"courseSelectionSystem/router"
	"github.com/gin-gonic/gin"
)

func main() {
	DB.RedisInit()
	DB.MysqlInit()
	DB.InitMemRedis()
	r := gin.Default()
	middle.InitSimpleSessionPool() //中间件 SimpleSessionPool
	r.Use(middle.HandleSimpleSession("camp-session"))
	router.RegisterRouter(r)
	r.Run(":8000")
}
