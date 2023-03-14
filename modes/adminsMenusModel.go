package modes

import (
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"strconv"
	"strings"
	"time"
)

type SysMenus struct {
	Id int `form:"id"`
	MenuName string `form:"menu_name"`
	MenuType int `form:"menu_type"`
	MenuOrder int `form:"menu_order"`
	MenuPid int `form:"menu_pid"`
	MenuUrl string `form:"menu_url"`
	MenuStatus int `form:"menu_status"`
	AddDatetime string `form:"add_datetime"`
}

type TreeList  struct {
	Id int `form:"id"`
	MenuName string `form:"menu_name"`
	MenuType int `form:"menu_type"`
	MenuOrder int `form:"menu_order"`
	MenuPid int `form:"menu_pid"`
	MenuUrl string `form:"menu_url"`
	MenuStatus int `form:"menu_status"`
	AddDatetime string `form:"add_datetime"`
	Children []*TreeList `form:"children"`
}

//搜索
func SearchMenuList(menu_name string) []*TreeList  {

	var menu  []SysMenus
	//获取第一级菜单
	DB.Model(SysMenus{}).Where("menu_name like ?","%"+menu_name+"%").Find(&menu)
	treeList := []*TreeList{}

	for  _,v := range(menu) {
		node := &TreeList{
			Id: v.Id,
			MenuName: v.MenuName,
			MenuType: v.MenuType,
			MenuOrder: v.MenuOrder,
			MenuPid: v.MenuPid,
			MenuUrl: v.MenuUrl,
			MenuStatus: v.MenuStatus,
			AddDatetime: v.AddDatetime,
		}
		node.Children = nil
		treeList = append(treeList,node)
	}

	return treeList
}
//menu单条
func MenuItem(id int) SysMenus  {
	var menu  SysMenus

	DB.Model(SysMenus{}).Where("id=?",id).Find(&menu)

	return menu
}

//menu单条
func CurrentMenu(c *gin.Context) SysMenus  {
	var menu  SysMenus

	var  url string
	//记录日志
	if(strings.Contains(c.Request.RequestURI,"?")){
		spurl := strings.Split(c.Request.RequestURI,"?")
		url = spurl[0]
	}else{
		url = c.Request.RequestURI
	}

	DB.Model(SysMenus{}).Where("menu_url=?",url).Find(&menu)

	return menu
}


//menu列表
func MenuList()  ([]*TreeList,int){
	return getMenu(0)
}
//获取菜单
func getMenu(pid int) ([]*TreeList,int) {
	var menu  []SysMenus
	var count = 0
	//获取第一级菜单
	DB.Model(SysMenus{}).Where("menu_pid=?",pid).Order("menu_order DESC ").Find(&menu)

	treeList := []*TreeList{}

	for  _,v := range(menu) {
		count++
		child,num:= getMenu(v.Id)
		node := &TreeList{
			Id: v.Id,
			MenuName: v.MenuName,
			MenuType: v.MenuType,
			MenuOrder: v.MenuOrder,
			MenuPid: v.MenuPid,
			MenuUrl: v.MenuUrl,
			MenuStatus: v.MenuStatus,
			AddDatetime: v.AddDatetime,
		}
		node.Children = child
		count += num
		treeList = append(treeList,node)
	}
	return treeList,count
}

//保存设置
func SaveMenu(menu *SysMenus) (bool,string) {

	if menu.Id == 0 {
		menu.AddDatetime = time.Now().Format(common.TimeTem)
		res := DB.Model(SysMenus{}).Create(menu)
		if res.Error != nil{
			return false,res.Error.Error()
		}
		return true,"保存成功"
	}else{
		res := DB.Model(SysMenus{}).Where("id = ?",menu.Id).Updates(menu)
		if res.Error != nil{
			return false,res.Error.Error()
		}
		return true,"保存成功"
	}
}

//保存设置
func DelMenu(id int) (bool,string) {
	var count int64
	DB.Model(SysMenus{}).Where("menu_pid = ?",id).Count(&count)
	if count>0 {
		return false,"菜单下有子菜单，删除失败"
	}

	res := DB.Delete(SysMenus{},id)
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"删除成功"
}

//权限判断
func HasPowerStr(admin_id int,menu_id int)  bool{
	//获取管理员对应的权限
	var admins = Admins{}
	var powers = AdminsPowers{}

	DB.Model(Admins{}).Where("admin_uid=?",admin_id).First(&admins)
	DB.Model(AdminsPowers{}).Where("admin_uid=?",admin_id).First(&powers)

	if admins.IsAdmin == 0 || powers.Id == 0  {
		return false
	}

	if admins.IsAdmin == 2 {
		return  true
	}
	id := strconv.Itoa(menu_id)
	ids := ","+id+","

	return strings.Contains(","+powers.AdminPowers+",",ids)
}

//权限判断
func HasPowerAction(admin_id int,menu_url string)  bool{
	//获取管理员对应的权限
	var admins = Admins{}
	var powers = AdminsPowers{}
	var menu = SysMenus{}

	DB.Model(Admins{}).Where("admin_uid=?",admin_id).First(&admins)
	DB.Model(AdminsPowers{}).Where("admin_uid=?",admin_id).First(&powers)
	DB.Model(SysMenus{}).Where("menu_url=?",menu_url).First(&menu)


	if admins.IsAdmin == 0 || powers.Id == 0 || menu.Id ==0  {
		return true
	}

	if admins.IsAdmin == 2 {
		return  false
	}
	id := strconv.Itoa(menu.Id)
	ids := ","+id+","

	return !strings.Contains(","+powers.AdminPowers+",",ids)
}


