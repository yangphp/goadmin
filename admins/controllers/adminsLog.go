package controllers

import (
	"github.com/gin-gonic/gin"
	"goadmin/modes"
	"net/http"
)

//日志列表 ，带翻页
func AdminLogList(c *gin.Context)  {
	//获取列表
	list := modes.AdminLogList(c)
	//获取管理员列表
	c.HTML(http.StatusOK,"admins/admin_log_list.html", gin.H{
		"list":list,
		"count":len(list),
		"start_date":c.Query("start_date"),
		"end_date":c.Query("end_date"),
		"admin_uname":c.Query("admin_uname"),
		"ip_addr":c.Query("ip_addr"),
	})
}

//日志详情
func AdminLogDetail(c *gin.Context)  {

}

//删除日志
func AdminLogDel(c *gin.Context)  {

	res,msg := modes.AdminLogDel(c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "删除成功",})

}

//批量删除日志
func AdminLogDelBatch(c *gin.Context)  {
	res,msg := modes.AdminLogDelBatch(c)
	if res != true {
		c.JSON(http.StatusOK, gin.H{ "code": 400, "msg":  msg,})
		return
	}
	c.JSON(http.StatusOK, gin.H{ "code": 0, "msg":  "删除成功",})
}



