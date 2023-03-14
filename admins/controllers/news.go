package controllers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"goadmin/modes"
	"net/http"
	"strconv"
)
//文章分类列表
func NewsCateList(c *gin.Context)  {

	list:=modes.NewsCateList(c)

	c.HTML(http.StatusOK,"news/news_cate_list.html", gin.H{
		"list":list,
		"cate_name":c.Query("cate_name"),
		"count":len(list),
	})
}
//文章分类添加/编辑
func NewsCateAdd(c *gin.Context)  {
	var cate modes.SysNewsCate
	news_id,exist := c.GetQuery("news_id");
	if  exist == true{
		//编辑
		news_idd,_ := strconv.Atoi(news_id)
		cate = modes.NewsCateItem(news_idd,c)
	}

	c.HTML(http.StatusOK,"news/news_cate_edit.html", gin.H{
		"cate":cate,
	})
}

//文章分类保存
func NewsCateSave(c *gin.Context)  {
	res,msg := modes.NewsCateSave(c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "保存成功",})
}
//文章分类删除
func NewsCateDel(c *gin.Context)  {

	res,msg := modes.NewsCateDel(c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "删除成功",})
}

//文章列表
func NewsList(c *gin.Context)  {

	list:=modes.NewsList(c,1)
	cate_list := modes.NewsCateList(c)

	//将cate_id 转换为int
	cate_id,_ := strconv.Atoi(c.Query("cate_id"))

	c.HTML(http.StatusOK,"news/news_list.html", gin.H{
		"list":list,
		"count":len(list),
		"cate_list":cate_list,
		"news_title":c.Query("news_title"),
		"cate_id":cate_id,
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
	})
}

//文章添加/编辑
func NewsAdd(c *gin.Context)  {

	news := modes.SysNews{}
	news_id,exist := c.GetQuery("news_id");
	if  exist == true{
		//编辑
		news_idd,_ := strconv.Atoi(news_id)
		news = modes.NewsItem(news_idd,c)
		//content base64解码
		bytes,_ := base64.StdEncoding.DecodeString(news.NewsContent)
		news.NewsContent = string(bytes)
	}
	cate_list := modes.NewsCateList(c)


	c.HTML(http.StatusOK,"news/news_edit.html", gin.H{
		"news":news,
		"cate_list":cate_list,
	})
}

//文章保存
func NewsSave(c *gin.Context)  {
	res,msg := modes.NewsSave(c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "保存成功",})
}

//文章删除
func NewsDel(c *gin.Context)  {

	news_id := c.Query("news_id")
	news_idd,_ := strconv.Atoi(news_id)
	res,msg := modes.NewsDel(news_idd,c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "保存成功",})
}


//文章列表
func NewsTrashList(c *gin.Context)  {
	list:=modes.NewsList(c,2)
	cate_list := modes.NewsCateList(c)

	//将cate_id 转换为int
	cate_id,_ := strconv.Atoi(c.Query("cate_id"))

	c.HTML(http.StatusOK,"news/news_trash_list.html", gin.H{
		"list":list,
		"count":len(list),
		"cate_list":cate_list,
		"news_title":c.Query("news_title"),
		"cate_id":cate_id,
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
	})
}

//恢复
func NewsRestore(c *gin.Context)  {
	news_id := c.Query("news_id")
	news_idd,_ := strconv.Atoi(news_id)
	res,msg := modes.NewsRestore(news_idd,c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "保存成功",})
}
//永久删除
func NewsDelReal(c *gin.Context)  {
	news_id := c.Query("news_id")
	news_idd,_ := strconv.Atoi(news_id)
	res,msg := modes.NewsDelReal(news_idd,c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "保存成功",})
}

