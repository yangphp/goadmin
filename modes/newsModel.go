package modes

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"strconv"
	"strings"
	"time"
)

type SysNewsCate struct {
	Id int `form:"id"`
	CateName string `form:"cate_name"`
	CateOrders string `form:"cate_orders"`
	IsShow int `form:"is_show"`
}

type SysNews struct {
	Id int `form:"id"`
	NewsTitle string `form:"news_title"`
	NewsDesc string `form:"news_desc"`
	NewsContent string `form:"news_content"`
	NewsPic string `form:"news_pic"`
	IsShow int `form:"is_show"`
	CateId int `form:"cate_id"`
	IsRecommed int `form:"is_recommed"`
	NewsReads int `form:"news_reads"`
	IsDel int `form:"is_del"`
	DelDatetime string `form:"del_datetime"`
	AddDatetime string `form:"add_datetime"`
}
type SysNewsList struct {
	SysNews
	CateName string
}

//获取新闻分类列表
func NewsCateList(c *gin.Context) []SysNewsCate{
	cate_list := []SysNewsCate{}

	NewDb := DB
	if cate_name,isExist := c.GetQuery("cate_name");isExist == true{
		if strings.TrimSpace(cate_name) != ""{
			NewDb = NewDb.Where("cate_name = ?",cate_name)
		}
	}
	res := NewDb.Model(SysNewsCate{}).Order("cate_orders DESC").Find(&cate_list)
	if res.Error != nil {
		return cate_list
	}
	return cate_list
}
//获取新闻分类
func NewsCateItem(news_id int,c *gin.Context) SysNewsCate {
	var news_cate  SysNewsCate
	DB.Model(SysNewsCate{}).Where("id = ?",news_id).First(&news_cate)

	return news_cate
}
//保存分类
func NewsCateSave(c *gin.Context) (bool,string) {
	var cate = SysNewsCate{}
	c.ShouldBind(&cate)

	if cate.Id > 0 {
		//更新
		res := DB.Model(SysNewsCate{}).Where("id=?",cate.Id).Updates(&cate)
		if res.Error != nil {
			return false,"保存失败"
		}

	}else{
		//增加
		res := DB.Model(SysNewsCate{}).Create(&cate)
		if res.Error != nil {
			return false,"保存失败"
		}
	}
	return true,"保存成功"
}
//删除分类
func NewsCateDel(c *gin.Context) (bool,string){
	id := c.Query("id")
	idd,err  := strconv.Atoi(id)
	if err !=nil {
		return false,err.Error()
	}
	var  count int64
	//判断分类下是否有文章
	 DB.Model(SysNews{}).Where("cate_id=?",idd).Count(&count)
	if count >0 {
		count_str := strconv.Itoa(int(count))
		return false,"该分类下有"+count_str+"篇文章，请先移除后再删除分类"
	}

	res :=DB.Model(SysNewsCate{}).Where("id=?",idd).Delete(SysNewsCate{})
	if res.Error != nil {
		return false,res.Error.Error()
	}
	return true,"保存成功"
}

//文章列表
func NewsList(c *gin.Context,flag int) ([]SysNewsList)  {

	news_list := []SysNewsList{}

	NewDb := DB
	if title,isExist := c.GetQuery("news_title");isExist == true{
		if strings.TrimSpace(title) != ""{
			NewDb = NewDb.Where("news_title like ?","%"+title+"%")
		}
	}
	if is_show,isExist := c.GetQuery("is_show");isExist == true{
		if strings.TrimSpace(is_show) != ""{
			NewDb = NewDb.Where("is_show = ?",is_show)
		}
	}

	if cate_id,isExist := c.GetQuery("cate_id");isExist == true{
		if strings.TrimSpace(cate_id) != ""{
			NewDb = NewDb.Where("cate_id = ?",cate_id)
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
	if flag == 1 {
		//只查询未删除的
		NewDb = NewDb.Where("is_del  = ?",1)
	}else{
		//只查询未删除的
		NewDb = NewDb.Where("is_del  = ?",2)
	}

	res := NewDb.Model(SysNews{}).Select("sys_news.*","sys_news_cate.cate_name").
		Order("sys_news.id DESC").
		Joins("join sys_news_cate on sys_news.cate_id = sys_news_cate.id").Find(&news_list)

	if res.Error != nil {
		return news_list
	}
	return news_list

}

//文章详情
func NewsItem(news_id int,c *gin.Context) SysNews  {
	var news  SysNews
	DB.Model(SysNews{}).Where("id = ?",news_id).First(&news)

	return news
}

//文章保存
func NewsSave(c *gin.Context) (bool,string) {
	var news = SysNews{}
	c.ShouldBind(&news)

	if news.Id > 0 {
		//更新

		//对内容进行base64编码一下
		content := []byte(news.NewsContent)
		news.NewsContent = base64.StdEncoding.EncodeToString(content)

		////解编码
		//bytes,_ := base64.StdEncoding.DecodeString(news.NewsContent)
		//news.NewsContent = string(bytes)

		res := DB.Model(SysNews{}).Where("id=?",news.Id).Updates(&news)
		if res.Error != nil {
			return false,"保存失败"
		}

	}else{
		//增加

		//对内容进行base64编码一下,方便保存
		content := []byte(news.NewsContent)
		news.NewsContent = base64.StdEncoding.EncodeToString(content)
		news.AddDatetime = time.Now().Format(common.TimeTem)
		news.DelDatetime = "1999-01-01 00:00:00"

		res := DB.Model(SysNews{}).Create(&news)
		if res.Error != nil {
			return false,"保存失败"
		}
	}
	return true,"保存成功"
}

//文章删除
func NewsDel(news_id int,c *gin.Context) (bool,string) {
	var news = SysNews{
		IsDel: 2,
		DelDatetime: time.Now().Format(common.TimeTem),
	}

	res := DB.Model(SysNews{}).Where("id=?",news_id).Updates(&news)
	if res.Error != nil {
		return false,"保存失败"
	}
	return true,"删除成功"
}

//文章恢复
func NewsRestore(news_id int,c *gin.Context) (bool,string) {
	var news = SysNews{
		IsDel: 1,
		DelDatetime: "1999-01-01 00:00:00",
	}

	res := DB.Model(SysNews{}).Where("id=?",news_id).Updates(&news)
	if res.Error != nil {
		return false,"保存失败"
	}
	return true,"删除成功"
}

//文章彻底删除
func NewsDelReal(news_id int,c *gin.Context) (bool,string) {

	var news = SysNews{}

	res := DB.Model(SysNews{}).Where("id=?",news_id).Delete(&news)
	if res.Error != nil {
		return false,"保存失败"
	}
	return true,"删除成功"
}

