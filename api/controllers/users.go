package controllers

import (
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"goadmin/modes"
	"net/http"
	"strconv"
)

// 定义结构体，接收表单提交的数据
type LoginUser struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//用户登录
func UserLogin(c *gin.Context)  {
	//用户登录
	var u LoginUser
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	//密码验证
	stat,msg,user := modes.UserLogin(u.Name,u.Password)
	if stat == false {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": msg})
		return
	}
	//登录成功，获取token返回给前端
	token := common.GetToken(u.Name,user.Uid)
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "login success",
		"token": token,
		"user_id":user.Uid,
	})
	return
}

//获取用户信息
func GetUserInfo(c *gin.Context)  {

	var claims11 = common.UserClaims{}

	claims,_ := c.Get("cliaims")
	claims11 = claims.(common.UserClaims)

	uid := claims11.UserId
	//username := claims11.UserName
	user := modes.UsersItem(uid)


	//username := claims.UserName
	//判断token
	c.JSON(http.StatusOK, gin.H{
		"code":   0,
		"msg":    "GetUserInfo Success",
		"claims": claims,
		"uid":user.Uid,
		"username":user.Username,
		"usertype":user.Utype,
	})
}
//退出
func UserLoginout(c *gin.Context)  {

	var claims11 = common.UserClaims{}

	claims,_ := c.Get("cliaims")
	claims11 = claims.(common.UserClaims)

	//将token设置为空值，并过期
	idstr := strconv.Itoa(claims11.UserId)
	common.SetRedisVal("token"+idstr,"",-1)

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "login out success",
	})
	return
}