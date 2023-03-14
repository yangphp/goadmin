package modes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"strings"
	"time"
)

type UsersSocre struct {
	Id int `form:"id"`
	Uid int `form:"uid"`
	BeginScore float64 `form:"begin_score"`
	ChangeScore float64 `form:"change_score"`
	EndScore float64 `form:"end_score"`
	ChangeInfo string `form:"change_info"`
	AddDatetime string `form:"add_datetime"`
}
type UsersScoreL struct {
	UsersSocre
	Username string `form:"username"`
	Mobile string `form:"mobile"`
	Utype int `form:"utype"`
}

//积分列表
func ScoreList(c *gin.Context) []UsersScoreL {
	users := []UsersScoreL{}

	NewDb := DB
	if username,isExist := c.GetQuery("username");isExist == true{
		if strings.TrimSpace(username) != ""{
			NewDb = NewDb.Where("users.username = ?",username)
		}
	}
	if uid,isExist := c.GetQuery("uid");isExist == true{
		if strings.TrimSpace(uid) != ""{
			NewDb = NewDb.Where("users_socre.uid = ?",uid)
		}
	}

	if mobile,isExist := c.GetQuery("mobile");isExist == true{
		if strings.TrimSpace(mobile) != ""{
			NewDb = NewDb.Where("users.mobile = ?",mobile)
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
	res := NewDb.Model(UsersSocre{}).Select("users_socre.*,users.username,users.mobile").
		Order("users_socre.id DESC").
		Joins("join users on users_socre.uid = users.uid").Find(&users)

	if res.Error != nil {
		return users
	}
	return users
}

//积分变更
func ScoreChange(c *gin.Context) (bool,string)  {
	var score = UsersSocre{}
	var user  = Users{}
	//获取变量
	c.ShouldBind(&score)

	fmt.Println(score)

	uid := score.Uid
	change_score := score.ChangeScore
	//开启事务
	tx := DB.Begin()
	tx.Model(Users{}).Where("uid=?",uid).First(&user)

	score.BeginScore = user.Uscore
	score.EndScore = user.Uscore+change_score
	score.AddDatetime = time.Now().Format(common.TimeTem)

	fmt.Println(score)
	//更新用户表
	user_row :=tx.Model(Users{}).Where("uid=?",uid).Updates(Users{
		Uscore: score.EndScore,
	}).RowsAffected
	score_row := tx.Model(UsersSocre{}).Create(&score).RowsAffected

	if user_row > 0 && score_row > 0 {
		tx.Commit()
		//添加日志
		AddUsersLog(uid,2,"管理员变更积分："+score.ChangeInfo,c)

		return true,"变更成功"
	}else{
		tx.Rollback()
		return false,"更新失败"
	}
}