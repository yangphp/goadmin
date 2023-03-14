package modes

import (
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"strconv"
	"strings"
	"time"
)

type UsersType struct {
	Id int `form:"id"`
	TypeName string `form:"type_name"`
	TypeInfo string `form:"type_info"`
}
type Users struct {
	Uid int `form:"uid"`
	Username string `form:"username"`
	Passwd string `form:"passwd"`
	Mobile string `form:"mobile"`
	Uface string `form:"uface"`
	Utype int `form:"utype"`
	Uscore float64 `form:"uscore"`
	Ustatus int `form:"ustatus"`
	Uremark string `form:"uremark"`
	LastDatetime string `form:"last_datetime"`
	AddDatetime string `form:"add_datetime"`
}
type UsersL struct {
	Users
	TypeName string `form:"type_name"`
}

type UsersLog struct {
	Id int `from:"id"`
	Uid int `from:"uid"`
	LogType int `from:"log_type"`
	LogInfo string `from:"log_info"`
	LogIp string `from:"log_ip"`
	AddDatetime string `from:"add_datetime"`
}
type UsersLogL struct {
	UsersLog
	Username string `form:"username"`
}
//角色列表
func UsersCateList(c *gin.Context) []UsersType {
	types := []UsersType{}
	DB.Model(UsersType{}).Find(&types)

	return types
}
//角色
func UsersCateItem(id int) UsersType {

	types := UsersType{}
	DB.Model(UsersType{}).Where("id=?",id).Find(&types)
	return types
}
//保存角色
func  UsersCateSave(c *gin.Context)  (bool,string){
	types := UsersType{}
	id,exist := c.GetPostForm("id")
	idd,_ := strconv.Atoi(id)
	c.ShouldBind(&types)
	if exist && idd >0{
		//修改角色

		res :=DB.Model(UsersType{}).Where("id=?",idd).Updates(&types)
		if res.Error != nil {
			return false,res.Error.Error()
		}else{
			return true,"修改成功"
		}
	}else{
		//新增角色
		res :=DB.Model(UsersType{}).Create(&types)
		if res.Error != nil {
			return false,res.Error.Error()
		}else{
			return true,"保存成功"
		}
	}
}

//删除角色
func  UsersCateDel(c *gin.Context) (bool,string) {
	id,exist := c.GetQuery("id")


	if exist {
		idd,_ := strconv.Atoi(id)
		//判断分类下是否有用户
		var count int64
		DB.Model(Users{}).Where("utype=?",idd).Count(&count)
		if count > 0 {
			count_str := strconv.Itoa(int(count))
			return false,"该分类下还有"+count_str+"用户，请先将分类下的用户移除后再删除分类"
		}else{
			res := DB.Where("id=?",idd).Delete(UsersType{})
			if res.Error!=nil {
				return false,res.Error.Error()
			}
			return true,"删除成功"
		}
	}
	return false,"删除失败"
}

//会员列表
func UsersList(c *gin.Context) []UsersL {
	users := []UsersL{}

	NewDb := DB
	if username,isExist := c.GetQuery("username");isExist == true{
		if strings.TrimSpace(username) != ""{
			NewDb = NewDb.Where("username like ?","%"+username+"%")
		}
	}
	if uid,isExist := c.GetQuery("uid");isExist == true{
		if strings.TrimSpace(uid) != ""{
			NewDb = NewDb.Where("uid = ?",uid)
		}
	}

	if mobile,isExist := c.GetQuery("mobile");isExist == true{
		if strings.TrimSpace(mobile) != ""{
			NewDb = NewDb.Where("mobile = ?",mobile)
		}
	}

	if utype,isExist := c.GetQuery("utype");isExist == true{
		if strings.TrimSpace(utype) != ""{
			NewDb = NewDb.Where("utype = ?",utype)
		}
	}

	if ustatus,isExist := c.GetQuery("ustatus");isExist == true{
		if strings.TrimSpace(ustatus) != ""{
			NewDb = NewDb.Where("ustatus = ?",ustatus)
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
	//连表查询
	res := NewDb.Model(Users{}).Select("users.*","users_type.type_name").
		Order("users.uid DESC").
		Joins("join users_type on users.utype = users_type.id").Find(&users)

	if res.Error != nil {
		return users
	}
	return users
}

//用户详情
func UsersItem(uid int) UsersL {
	users := UsersL{}
	DB.Model(Users{}).Select("users.*","users_type.type_name").
		Joins("join users_type on users.utype = users_type.id").
		Where("users.uid=?",uid).Find(&users)

	return users
}
//用户详情
func UsersItemByName(username string) (bool,Users) {
	users := Users{}
	res := DB.Model(Users{}).Where("username=?",username).Find(&users)
	if res.Error != nil {
		return false,users
	}
	if users.Uid >0 {
		return true,users
	}
	return false,users
}
//保存用户信息
func  UsersSave(c *gin.Context)  (bool,string){
	users := Users{}
	id,exist := c.GetPostForm("uid")
	idd,_ := strconv.Atoi(id)
	c.ShouldBind(&users)

	users.LastDatetime = time.Now().Format(common.TimeTem)

	if exist && idd >0{
		//修改用户
		if strings.TrimSpace(users.Passwd)!="" {
			users.Passwd = common.MyMd5(users.Passwd)
		}
		res :=DB.Model(Users{}).Where("uid=?",idd).Updates(&users)
		if res.Error != nil {
			return false,res.Error.Error()
		}else{
			return true,"修改成功"
		}
	}else{
		//新增会员
		users.AddDatetime = time.Now().Format(common.TimeTem)
		users.Passwd = common.MyMd5(users.Passwd)

		res :=DB.Model(Users{}).Create(&users)
		if res.Error != nil {
			return false,res.Error.Error()
		}else{
			return true,"保存成功"
		}
	}
}

//删除角色
func  UsersDel(c *gin.Context) (bool,string) {
	id,exist := c.GetQuery("uid")

	if exist {
		idd,_ := strconv.Atoi(id)

		res := DB.Where("uid=?",idd).Delete(Users{})
		if res.Error!=nil {
			return false,res.Error.Error()
		}
		return true,"删除成功"
	}
	return false,"删除失败"
}

//添加日志
// log_type 1登录退出，2积分变更 3其他提示
func AddUsersLog(uid int,log_type int,log_info string,c *gin.Context) (bool,string) {
	var logs = UsersLog{
		Uid:uid,
		LogType: log_type,
		LogInfo: log_info,
		LogIp: c.RemoteIP(),
		AddDatetime: time.Now().Format(common.TimeTem),
	}

	DB.Model(UsersLog{}).Create(&logs)
	return true,"添加日志成功"
}

//日志列表
func UsersLogList(c *gin.Context) []UsersLogL {

	logl := []UsersLogL{}
	NewDb := DB
	if username,isExist := c.GetQuery("username");isExist == true{
		if strings.TrimSpace(username) != ""{
			NewDb = NewDb.Where("users.username like ?","%"+username+"%")
		}
	}
	if uid,isExist := c.GetQuery("uid");isExist == true{
		if strings.TrimSpace(uid) != ""{
			NewDb = NewDb.Where("users_log.uid = ?",uid)
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
	//连表查询
	res := NewDb.Model(UsersLog{}).Select("users_log.*,users.username").
		Order("users_log.id DESC").
		Joins("join users on users.uid = users_log.uid").Find(&logl)

	if res.Error != nil {
		return logl
	}
	return logl
}

//删除日志
func  UsersLogDel(c *gin.Context) (bool,string) {

	id,exist := c.GetQuery("id")

	if exist {
		idd,_ := strconv.Atoi(id)

		res := DB.Where("id=?",idd).Delete(UsersLog{})
		if res.Error!=nil {
			return false,res.Error.Error()
		}
		return true,"删除成功"
	}
	return false,"删除失败"
}

//用户登录
func UserLogin(name string,pass string) (bool,string,Users) {

	var user = Users{}
	pass = common.MyMd5(pass)

	DB.Model(Users{}).Where("username=? and passwd = ?",name,pass).First(&user)
	if  user.Uid == 0{
		return false,"登录失败",user
	}
	if user.Ustatus == 2 {
		return false,"用户处于风控中，暂无法登录",user
	}
	if user.Ustatus == 2 {
		return false,"用户已被封禁，暂无法登录",user
	}
	return  true,"登录成功",user
}
