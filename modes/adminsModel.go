package modes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"net/http"
	"strings"
	"time"
)

var DB = common.DB
//管理员
type Admins struct {
	AdminUid int `form:"admin_uid"`
	AdminUname string `form:"admin_uname"`
	AdminPasswd string `form:"admin_passwd"`
	AdminTruename string `form:"admin_truename"`
	AdminDept string `form:"admin_dept"`
	Lastlogin string `form:"last_login"`
	Logincount int `form:"logincount"`
	IsAdmin int `form:"is_admin"`
	AdminStatus int `form:"admin_status"`
	AdminRemark string `form:"admin_remark"`
	AddDatetime string `form:"add_datetime"`
}
//权限表
type AdminsPowers struct {
	Id int `form:"id"`
	AdminUid int `form:"auid"`
	AdminPowers string `form:"admin_powers"`
	LastDatetime string `form:"last_datetime"`
	AddDatetime string `form:"add_datetime"`
}
//日志记录表
type AdminsLog struct {
	Id int `form:"id"`
	AdminUid int `form:"admin_uid"`
	RequestUrl string `form:"request_url"`
	RequestMethod string `form:"request_method"`
	RequestParams string `form:"request_params"`
	IpAddr string	`form:"ip_addr"`
	AddDatetime string `form:"add_datetime"`
}

//管理员注册
func AdminsAdd(admin Admins,c *gin.Context) (bool,string) {
	//密码加密
	admin.AdminPasswd = common.MyMd5(admin.AdminPasswd)
	admin.AddDatetime = time.Now().Format(common.TimeTem)
	admin.Lastlogin = time.Now().Format(common.TimeTem)

	//判断用户名是否重复
	var count int64
	DB.Model(Admins{}).Where("admin_uname=?",admin.AdminUname).Count(&count)
	if count > 0 {
		return false,"该管理员账号已存在，请重新设置其他账号"
	}

	res := DB.Model(Admins{}).Create(&admin)
	if res.RowsAffected == 1{
		return true,"创建成功"
	}
	return false,res.Error.Error()
}
//管理员修改
func AdminsSave(admin Admins,c *gin.Context) (bool,string) {

	//密码加密
	if strings.TrimSpace(admin.AdminPasswd) != ""{
		admin.AdminPasswd = common.MyMd5(admin.AdminPasswd)
	}

	DB.Model(Admins{}).Where("admin_uid=?",admin.AdminUid).Updates(&admin)
	return true,"创建成功"
}
//管理员修改密码
func AdminsChangePass(admin_id int,oldpass string,newpass string,) (bool,string) {
	var admin = Admins{}
	//密码加密
	if strings.TrimSpace(oldpass) != ""{
		oldpass = common.MyMd5(oldpass)
	}

	if strings.TrimSpace(newpass) != ""{
		newpass = common.MyMd5(newpass)
	}
	DB.Model(Admins{}).Where("admin_uid=?",admin_id).First(&admin)
	if admin.AdminPasswd != oldpass {
		return false,"旧密码输入错误"
	}

	DB.Model(Admins{}).Where("admin_uid=?",admin.AdminUid).Updates(Admins{
		AdminPasswd: newpass,
	})
	return true,"密码修改成功"
}
//获取管理员
func AdminsItem(aid int) (Admins) {

	var info Admins
	DB.Model(Admins{}).Where("admin_uid=?",aid).Find(&info)

	return info
}

//获取管理员账号
func AdminsUname(aid int) (string) {

	//var uname string
	var users Admins
	DB.Model(Admins{}).Select("admin_uname").Where("admin_uid=?",aid).First(&users)

	return users.AdminUname
}

//删除管理员
func AdminDel(uid int,c *gin.Context) (bool,string) {

	//goadmin uid=6的不能被删除
	if uid == 6 {
		return false,"系统超级管理员不能被删除！"
	}
	//不能自己删除自己
	admin_id,exist := c.Get("admin_uid")
	if exist && admin_id == uid {
		return false,"不能删除当前登录的账号！"
	}

	//判断权限
	res := DB.Delete(Admins{},uid)
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"删除成功"
}
//更改状态
func AdminStatus(uid int,status string) (bool,string) {

	//goadmin uid=6的不能被删除
	if uid == 6 {
		return false,"系统超级管理员不能被停用！"
	}
	//判断权限
	res := DB.Model(Admins{}).Where("admin_uid=?",uid).Update("admin_status",status)
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"删除成功"
}

//管理员列表
func AdminList(c *gin.Context) ([]Admins) {
	users := []Admins{}
	NewDb := DB
	if admin_uname,isExist := c.GetQuery("admin_uname");isExist == true{
		if strings.TrimSpace(admin_uname) != ""{
			NewDb = NewDb.Where("admin_uname = ?",admin_uname)
		}

	}
	if admin_truename,isExist := c.GetQuery("admin_truename");isExist == true{
		if strings.TrimSpace(admin_truename) != ""{
			NewDb = NewDb.Where("admin_truename = ?",admin_truename)
		}
	}
	if start_date,isExist := c.GetQuery("start_date");isExist == true{
		if strings.TrimSpace(start_date) != "" {
			start_date = start_date + " 00:00:00"
			NewDb = NewDb.Where("add_datetime > ?", start_date)
		}
	}
	if end_date,isExist := c.GetQuery("end_date");isExist == true{
		if strings.TrimSpace(end_date) != "" {
			end_date = end_date+" 23:59:59"
			NewDb = NewDb.Where("add_datetime <= ?",end_date)
		}
	}

	res := NewDb.Model(Admins{}).Find(&users)
	if res.Error != nil {
		return users
	}
	return users
}
//管理员登录
func AdminLogin(name string,pass string,c *gin.Context) (common.ReData) {
	admin := Admins{}
	pass = common.MyMd5(pass)

	res := DB.Table("admins").Where("admin_uname=?",name).Where("admin_passwd=?",pass).
		Find(&admin)

	if res.Error != nil {
		return common.ReData{false,res.Error.Error(),&admin,}
	}
	//查询结果失败
	if admin.AdminUid == 0 {
		return common.ReData{false,"账号或密码输入错误",admin,}
	}

	lastlogin:=time.Now().Format(common.TimeTem)
	//更新登录次数，和最后登录时间
	upres := DB.Model(&admin).Where("admin_uid=?",admin.AdminUid).
		Updates(Admins{Lastlogin: lastlogin,Logincount:admin.Logincount+1})

	if upres.Error != nil {
		return common.ReData{false,res.Error.Error(),admin,}
	}
	//写入日志
	AddAdminLog(admin.AdminUid,c)
	return common.ReData{true,"登录成功",admin,}
}

//获取管理员权限
func PowerItem(aid int) (AdminsPowers) {

	var p AdminsPowers
	DB.Model(AdminsPowers{}).Where("admin_uid=?",aid).Find(&p)

	return p
}

//保存管理员权限
func PowerSave(aaid int,is_admin int,power string) (bool) {
	//更新权限状态
	DB.Model(Admins{}).Where("admin_uid=?",aaid).Update("is_admin",is_admin)

	var num int64
	DB.Model(AdminsPowers{}).Where("admin_uid=?",aaid).Count(&num)
	if num > 0 {
		//修改
		DB.Model(AdminsPowers{}).Where("admin_uid=?",aaid).Updates(AdminsPowers{
			AdminPowers:power,
			LastDatetime: time.Now().Format(common.TimeTem),
		})
	}else{
		DB.Model(AdminsPowers{}).Create(&AdminsPowers{
			AdminUid: aaid,
			AddDatetime: time.Now().Format(common.TimeTem),
			AdminPowers:power,
			LastDatetime: time.Now().Format(common.TimeTem),
		})
	}
	return true
}

//获取管理员日志列表
func AdminLogList(c *gin.Context) ([]AdminsLog) {
	log_list := []AdminsLog{}
	admin :=Admins{}
	NewDb := DB

	if admin_uname,isExist := c.GetQuery("admin_uname");isExist == true{
		if strings.TrimSpace(admin_uname) != ""{
			res := DB.Model(Admins{}).Where("admin_uname=?",admin_uname).Find(&admin)
			if res.Error == nil {
				NewDb = NewDb.Where("admin_uid = ?",admin.AdminUid)
			}
		}
	}
	if ip_addr,isExist := c.GetQuery("ip_addr");isExist == true{
		if strings.TrimSpace(ip_addr) != ""{
			if strings.Contains(ip_addr,":") {
				ipa := strings.Split(ip_addr,":")
				ip_addr = ipa[0]
			}
			NewDb = NewDb.Where("ip_addr like ?",ip_addr+"%")
		}
	}
	if start_date,isExist := c.GetQuery("start_date");isExist == true{
		if strings.TrimSpace(start_date) != "" {
			start_date = start_date + " 00:00:00"
			NewDb = NewDb.Where("add_datetime > ?", start_date)
		}
	}
	if end_date,isExist := c.GetQuery("end_date");isExist == true{
		if strings.TrimSpace(end_date) != "" {
			end_date = end_date+" 23:59:59"
			NewDb = NewDb.Where("add_datetime <= ?",end_date)
		}
	}

	res := NewDb.Model(AdminsLog{}).Find(&log_list)
	if res.Error != nil {
		return log_list
	}
	return log_list
}

//写入登录日志
func AddAdminLog(admin_uid int,c *gin.Context) bool {
	var  url string
	var params_str string
	//记录日志
	if(strings.Contains(c.Request.RequestURI,"?")){
		spurl := strings.Split(c.Request.RequestURI,"?")
		url = spurl[0]
	}else{
		url = c.Request.RequestURI
	}

	if (c.Request.Method == "POST") {

		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				fmt.Println(nil, err)
			}
		}

		var postMap = make(map[string]any, len(c.Request.PostForm))
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
		//密码，加密存储
		if url == "/admin/dologin" {
			postMap["passwd"]  = common.MyMd5(postMap["passwd"].(string))
		}

		dataType , _ := json.Marshal(postMap)
		params_str = string(dataType)

	}else{
		query := c.Request.URL.Query()
		var queryMap = make(map[string]any, len(query))
		for k := range query {
			queryMap[k] = c.Query(k)
		}
		dataType1 , _ := json.Marshal(queryMap)
		params_str = string(dataType1)
	}

	admins_log := AdminsLog{
		AdminUid: admin_uid,
		RequestUrl: url,
		RequestMethod: c.Request.Method,
		RequestParams: params_str,
		IpAddr: c.Request.RemoteAddr,
		AddDatetime: time.Now().Format(common.TimeTem),
	}
	res := DB.Model(AdminsLog{}).Create(&admins_log)

	if res.RowsAffected == 1{
		return true
	}
	return false
}

//删除日志
func AdminLogDel(c *gin.Context) (bool,string) {
	//判断权限
	lid := c.PostForm("id")
	res := DB.Delete(AdminsLog{},lid)
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"删除成功"
}
func AdminLogDelBatch(c *gin.Context) (bool,string) {
	//判断权限
	var ids []string
	lid := c.PostForm("ids")
	ids = strings.Split(lid,",")

	res := DB.Model(AdminsLog{}).Where(" id in ?", ids).Delete(&AdminsLog{})
	if res.Error != nil{
		return false,res.Error.Error()
	}
	return true,"删除成功"
}
