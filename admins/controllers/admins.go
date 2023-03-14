package controllers

import (
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"goadmin/modes"
	"net/http"
	"strconv"
	"strings"
)

var db = common.DB

//后台首页
func AdminIndex(c *gin.Context)  {
	//会员中心主页
	admin_uid,_ :=c.Get("admin_uid")

	//获取总共登录次数
	var count int64
	db.Table("admins_log").Where("admin_uid=? and request_url= '/admin/dologin'",admin_uid).Count(&count)
	//获取上次登录IP和上次登录时间
	var admins_log = modes.AdminsLog{}
	db.Table("admins_log").Where("admin_uid=? and request_url= '/admin/dologin' ",admin_uid).Order("id Desc").First(&admins_log)


	redata := gin.H{
		"count":count,
		"admins_log":&admins_log,
	}
	c.HTML(http.StatusOK,"admins/index.html", redata)
}

//管理员退出
func LoginOut(c *gin.Context)  {
	common.ClearCookie("admin_uid",c)
	common.ClearCookie("admin_auth",c)
	common.ClearCookie("sessionid",c)

	common.ClearSession("admin_uid",c)
	common.ClearSession("admin_auth",c)
	common.ClearSession("sessionid",c)

	c.Redirect(http.StatusFound,"/admin/login")
	c.Abort() //不加这个abort后续代码还会执行
	return
}

//管理员列表
func AdminList(c *gin.Context)  {
	//获取列表
	users := modes.AdminList(c)
	//获取管理员列表
	c.HTML(http.StatusOK,"admins/list.html", gin.H{
		"list":users,
		"count":len(users),
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
		"admin_uname":c.Query("admin_uname"),
		"admin_truename":c.Query("admin_truename"),
	})
}

//管理员添加
func AdminAdd(c *gin.Context)  {
	//获取管理员列表
	c.HTML(http.StatusOK,"admins/admin_add.html", gin.H{})
}



//管理员编辑
func AdminEdit(c *gin.Context)  {
	aidd,_ := strconv.Atoi(c.Query("aid"))
	var info modes.Admins
	info = modes.AdminsItem(aidd)

	//获取管理员列表
	c.HTML(http.StatusOK,"admins/admin_edit.html", gin.H{
		"info":info,
	})
}

//管理员详情
func AdminInfo(c *gin.Context)  {
	aidd,_ := strconv.Atoi(c.Query("admin_id"))
	var info modes.Admins
	info = modes.AdminsItem(aidd)

	//获取管理员列表
	c.HTML(http.StatusOK,"admins/admin_info.html", gin.H{
		"info":info,
	})
}
//修改密码
func AdminChangePass(c *gin.Context)  {
	aidd,_ := strconv.Atoi(c.Query("admin_id"))
	var info modes.Admins
	info = modes.AdminsItem(aidd)

	//获取管理员列表
	c.HTML(http.StatusOK,"admins/admin_change_pass.html", gin.H{
		"info":info,
	})
}
//保存密码
func AdminSavePass(c *gin.Context)  {

	admin_id,_ := strconv.Atoi(c.PostForm("admin_id"))
	old_pass,_ := c.GetPostForm("old_pass")
	new_pass,_ := c.GetPostForm("new_pass")
	re_pass,_ := c.GetPostForm("re_pass")

	if new_pass != re_pass || new_pass == "" {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  "新密码两次输入不一致或为空",})
		return
	}

	res,msg :=modes.AdminsChangePass(admin_id,old_pass,new_pass)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "密码设置成功，请重新登录",})

}
//管理员保存
func AdminSave(c *gin.Context)  {

	admin := modes.Admins{}
	c.ShouldBind(&admin)

	if admin.AdminUid > 0 {
		res,msg := modes.AdminsSave(admin,c)
		if res != true {
			c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
			return
		}
		c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "修改成功",})

	}else{
		res,msg:= modes.AdminsAdd(admin,c)
		if res != true {
			c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
			return
		}
		c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "添加成功",})
	}
}

//管理员保存
func AdminDel(c *gin.Context)  {

	uid := c.Query("uid")
	uidd,_ := strconv.Atoi(uid)
	stat,msg := modes.AdminDel(uidd,c)

	if stat != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "删除成功",})
}

//更改状态
func  AdminStatus(c *gin.Context){
	uid := c.Query("uid")
	status := c.Query("admin_status")
	uidd,_ := strconv.Atoi(uid)
	stat,msg := modes.AdminStatus(uidd,status)

	if stat != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "修改成功",})
}

//管理员权限
func AdminPower(c *gin.Context)  {
	aidd,_ := strconv.Atoi(c.Query("aid"))
	var info modes.Admins
	info = modes.AdminsItem(aidd)
	//获取权限
	power := modes.PowerItem(aidd)
	power_str := ","+power.AdminPowers+","
	//获取菜单列表
	menus,_ := modes.MenuList()
	//获取管理员权限
	c.HTML(http.StatusOK,"admins/admin_power.html", gin.H{
		"info":info,
		"power":power,
		"menus":menus,
		"power_str":power_str,
	})
}

//管理员权限
func AdminPowerSave(c *gin.Context)  {
	//获取所有参数
	aaid,_ := strconv.Atoi(c.PostForm("aaid"))
	admin_powers :=c.PostFormArray("pids[]")
	is_admin,_ := strconv.Atoi(c.PostForm("is_admin"))

	p := strings.Join(admin_powers,",")
	if p == "" {
		p = "0,"
	}

	res := modes.PowerSave(aaid,is_admin,p)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  "设置权限失败",})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "设置权限成功",})
}