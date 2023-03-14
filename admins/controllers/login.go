package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"goadmin/common"
	"goadmin/modes"
	"net/http"
	"os"
	"strconv"
)

type LoginForm struct {
	UserName string `form:"username" binding:"required,alphanum,min=6,max=20"`
	Passwd string `form:"passwd" bindding:"required,alphanum,min=6,max=20"`
	Captcha   string `form:"captcha" binding:"required,capt"`
	CaptchaId string `form:"captcha_id" bingding:"required"`
	Online string `form:"online"`
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

//登录页面
func AdminLogin(c *gin.Context)  {
	//显示登录模板
	c.HTML(http.StatusOK,"login/index.html",nil)
}
//登录处理
func AdminLoginAction(c *gin.Context)  {

	var u LoginForm
	err := c.ShouldBind(&u)
	if err != nil {
		// 401 验证码错误
		c.JSON(http.StatusOK, gin.H{ "code": 401, "msg":  err.Error(),})
		return
	}
	//验证账号密码
	res := modes.AdminLogin(u.UserName,u.Passwd,c)
	if res.Status != true {
		c.JSON(http.StatusOK, gin.H{ "code": 402, "msg":  res.Msg,})
		return
	}
	//类型断言
	admins,_ := (res.Data).(modes.Admins)

	//session过期时间为2小时
	common.SetSession("admin_uid",strconv.Itoa(admins.AdminUid),c)
	common.SetSession("admin_auth",common.MyMd5(strconv.Itoa(admins.AdminUid)+admins.AdminUname),c)

	if u.Online == "1" {
		//保持登录状态1周 记录session 和cookie
		common.SetCookie("admin_uid",strconv.Itoa(admins.AdminUid),c)
		common.SetCookie("admin_auth",""+common.MyMd5(strconv.Itoa(admins.AdminUid)+admins.AdminUname),c)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
	return
}


//无路由
func AdminError(c *gin.Context)  {
	//显示登录模板
	c.HTML(http.StatusOK,"login/404.html",nil)
}


