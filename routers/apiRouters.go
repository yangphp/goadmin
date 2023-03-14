package routers

import (
	"github.com/gin-gonic/gin"
	"goadmin/api/controllers"
	"goadmin/api/middleware"
)

func ApiRouter(router *gin.Engine)  {

	//会员登录
	router.POST("users/login",controllers.UserLogin)

	//使用JWT对用户的请求进行验证
	user:=router.Group("users/", middleware.CheckAuth)
	{
		//获取会员信息
		user.POST("getuserinfo",controllers.GetUserInfo)
		user.POST("loginout",controllers.UserLoginout)
	}
}

