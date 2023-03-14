package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"goadmin/modes"
	"net/http"
	"strconv"
	"strings"
)
type UserInfo struct {
	UserName string
	UserTrueName string
	AdminUid int
}

var userinfo = &UserInfo{"","",0}
var menus []*modes.TreeList
var current_menu modes.SysMenus

func CheckLogin (c *gin.Context)  {

	admin_uid2  := common.GetCookie("admin_uid",c)
	admin_auth2  := common.GetCookie("admin_auth",c)

	var admin_uid string
	var admin_auth string


	if admin_uid2  != "" && admin_auth2  != ""{
		admin_uid = admin_uid2
		admin_auth = admin_auth2
	}else{
		//重定向到登录页面
		c.Redirect(http.StatusFound,"/admin/login")
		c.Abort() //不加这个abort后续代码还会执行
		return
	}

	//判断用户名密码，是否正常
	myDb := common.DB
	admins := modes.Admins{}
	//获取admin
	myDb.Model(&modes.Admins{}).First(&admins,admin_uid)
	fmt.Println(admins)
	check_token := common.MyMd5(strconv.Itoa(admins.AdminUid)+admins.AdminUname)


	if check_token !=  admin_auth || admin_auth == ""{
		fmt.Println("====1111111======")
		//重定向到登录页面
		c.Redirect(http.StatusFound,"/admin/login")
		c.Abort() //不加这个abort后续代码还会执行
		return
	}

	//权限判断
	var  url string
	var menu  = modes.SysMenus{}
	//记录日志
	if(strings.Contains(c.Request.RequestURI,"?")){
		spurl := strings.Split(c.Request.RequestURI,"?")
		url = spurl[0]
	}else{
		url = c.Request.RequestURI
	}
	myDb.Model(modes.SysMenus{}).Where("menu_url=?",url).First(&menu)
	if menu.Id>0 {
		res :=modes.HasPowerStr(admins.AdminUid,menu.Id)
		if res != true {
			c.JSON(http.StatusOK,gin.H{
				"code":400,
				"msg":"无此操作权限",
			})
			c.Abort() //不加这个abort后续代码还会执行
			return
		}
	}

	//设置值
	c.Set("admin_uid", admins.AdminUid)
	c.Set("admin_info",admins)

	userinfo.AdminUid = admins.AdminUid
	userinfo.UserName = admins.AdminUname
	userinfo.UserTrueName = admins.AdminTruename
	//记录日志
	modes.AddAdminLog(admins.AdminUid,c)
	//获取菜单列表
	menus,_ = modes.MenuList()
	//获取当前的菜单
	current_menu = modes.CurrentMenu(c)

}

func GetAdminName() string  {
	return userinfo.UserName
}
func GetAdminId() int  {
	return userinfo.AdminUid
}
func GetAdminTruename() string  {
	return userinfo.UserTrueName
}

func GetMenus() []*modes.TreeList {
	return menus
}
func GetCurrentMenu() modes.SysMenus {
	return current_menu
}
