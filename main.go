package main

import (
	"github.com/gin-gonic/gin"
	"goadmin/admins/controllers"
	"goadmin/routers"
)


func main()  {
	//使用gin框架 1.9
	router := gin.Default()
	//加载管理后台路由
	routers.AdminRouter(router)
	//加载API路由
	routers.ApiRouter(router)

	//没有匹配上路由
	router.NoRoute(controllers.AdminError)
	//运行到8088端口
	router.Run(":8088")
}
