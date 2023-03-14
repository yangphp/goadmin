package controllers

import (
	"github.com/gin-gonic/gin"
	"goadmin/modes"
	"net/http"
	"strconv"
)

//角色列表
func UsersCateList(c *gin.Context)  {

	list := modes.UsersCateList(c)

	c.HTML(http.StatusOK,"users/users_cate_list.html",gin.H{
		"list":list,
		"count":len(list),
	})
}

//添加，编辑角色
func UsersCateAdd(c *gin.Context)  {

	types := modes.UsersType{}
	id,exist := c.GetQuery("id")
	if exist {
		idd,_ := strconv.Atoi(id)
		types = modes.UsersCateItem(idd)
	}

	c.HTML(http.StatusOK,"users/users_cate_add.html",gin.H{
		"cate":types,
	})
}

//保存角色
func UsersCateSave(c *gin.Context)  {
	res,msg := modes.UsersCateSave(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}
}

//删除角色
func UsersCateDel(c *gin.Context)  {
	res,msg := modes.UsersCateDel(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}
}

//会员列表
func UsersList(c *gin.Context)  {
	list := modes.UsersList(c)
	type_list := modes.UsersCateList(c)

	utype := c.Query("utype")
	utype_id,_ := strconv.Atoi(utype)

	c.HTML(http.StatusOK,"users/users_list.html",gin.H{
		"list":list,
		"count":len(list),
		"type_list":type_list,
		"utype":utype_id,
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
		"uid":c.Query("uid"),
		"username":c.Query("username"),
		"mobile":c.Query("mobile"),
	})
}

//添加角色
func UsersAdd(c *gin.Context)  {

	type_list := modes.UsersCateList(c)
	c.HTML(http.StatusOK,"users/users_add.html",gin.H{
		"type_list":type_list,
	})
}
//查看会员
func UsersInfo(c *gin.Context)  {
	var users modes.UsersL

	uid,_ := c.GetQuery("uid")
	uidd,_ := strconv.Atoi(uid)

	users = modes.UsersItem(uidd)

	c.HTML(http.StatusOK,"users/users_info.html",gin.H{
		"users":users,
	})
}
//编辑角色
func UsersEdit(c *gin.Context)  {
	var users modes.UsersL

	uid,_ := c.GetQuery("uid")
	uidd,_ := strconv.Atoi(uid)

	users = modes.UsersItem(uidd)
	type_list := modes.UsersCateList(c)

	c.HTML(http.StatusOK,"users/users_add.html",gin.H{
		"type_list":type_list,
		"users":users,
	})
}

//保存会员
func UsersSave(c *gin.Context)  {
	res,msg := modes.UsersSave(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}
}

//删除会员
func UsersDel(c *gin.Context)  {
	res,msg := modes.UsersDel(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}
}
//积分列表
func SocreList(c *gin.Context)  {
	list := modes.ScoreList(c)

	c.HTML(http.StatusOK,"users/users_score_list.html",gin.H{
		"list":list,
		"count":len(list),
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
		"uid":c.Query("uid"),
		"username":c.Query("username"),
		"mobile":c.Query("mobile"),
	})
}

//变更积分
func SocreChange(c *gin.Context)  {

	users := modes.UsersL{}
	id,exist := c.GetQuery("uid")
	if exist {
		idd,_ := strconv.Atoi(id)
		users = modes.UsersItem(idd)
	}

	c.HTML(http.StatusOK,"users/users_score_change.html",gin.H{
		"users":users,
	})
}
//获取用户积分
func  GetUserScore(c *gin.Context){

	username,exist := c.GetQuery("username")
	if exist {
		stat,users := modes.UsersItemByName(username)
		if stat != true {
			c.JSON(http.StatusOK,gin.H{
				"code":400,
			})
			return
		}else{
			c.JSON(http.StatusOK,gin.H{
				"code":0,
				"uid":users.Uid,
				"score":users.Uscore,
			})
			return
		}
	}
	c.JSON(http.StatusOK,gin.H{"code":400,})
	return
}
//保存积分
func SocreSave(c *gin.Context)  {

	res,msg := modes.ScoreChange(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}

}

//会员日志列表
func  UsersLogList(c *gin.Context)  {
	list :=modes.UsersLogList(c)

	c.HTML(http.StatusOK,"users/users_log_list.html",gin.H{
		"list":list,
		"count":len(list),
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
		"uid":c.Query("uid"),
		"username":c.Query("username"),
	})
}
//删除日志
func UsersLogDel(c *gin.Context){
	res,msg := modes.UsersLogDel(c)
	if res!=true {
		c.JSON(http.StatusOK,gin.H{"code":400,"msg":msg})
	}else{
		c.JSON(http.StatusOK,gin.H{"code":0,"msg":msg})
	}
}