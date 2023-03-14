package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goadmin/modes"
	"net/http"
	"strconv"
)

//系统设置
func Setting(c *gin.Context)  {
	//获取系统设置
	setting := modes.GetSetting()

	//显示登录模板
	c.HTML(http.StatusOK,"system/setting.html",gin.H{
		"code"	: 0,
		"setting": setting,
	})
}
//保存系统设置
func DoSetting(c *gin.Context)  {
	var setting  modes.SysSetting
	err := c.ShouldBind(&setting)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	res,msg := modes.SaveSetting(&setting)
	if res != true {
		c.String(http.StatusOK, msg)
		return
	}
	//重定向到设置页面
	c.Redirect(http.StatusFound,"/admin/system_setting")
	return
}

//菜单列表
func Menus(c *gin.Context)  {
	//搜索
	menu_name := c.Query("menu_name")
	if menu_name != "" {
		//获取系统设置
		treeList := modes.SearchMenuList(menu_name)
		//显示登录模板
		c.HTML(http.StatusOK,"system/menu.html",gin.H{
			"treeList":treeList,
			"sum":len(treeList),
			"menu_name":menu_name,
		})
	}else{
		//获取系统设置
		treeList,sum := modes.MenuList()
		//显示登录模板
		c.HTML(http.StatusOK,"system/menu.html",gin.H{
			"treeList":treeList,
			"sum":sum,
			"menu_name":"",
		})
	}


}

//新增/编辑菜单
func MenusEdit(c *gin.Context)  {

	//获取分类
	treeList,_ := modes.MenuList()
	menu := modes.SysMenus{}
	//获取分类
	id,_ := strconv.Atoi(c.Query("id"))
	fmt.Println(id)
	if id != 0  {
		menu = modes.MenuItem(id)
	}
	//显示登录模板
	c.HTML(http.StatusOK,"system/menu_edit.html",gin.H{
		"treeList": treeList,
		"menu":menu,
	})
}

//保存菜单
func MenusSave(c *gin.Context)  {

	var menu  modes.SysMenus
	err := c.ShouldBind(&menu)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  err.Error(),})
		return
	}
	res,msg := modes.SaveMenu(&menu)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "添加成功",})
}

//删除菜单
func MenusDel(c *gin.Context)  {
	//获取系统设置
	id := c.PostForm("id")

	ids,_ :=strconv.Atoi(id)
	res,msg := modes.DelMenu(ids)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  msg,})
}