package common

import "github.com/gin-gonic/gin"

func SetCookie(cname string ,cvalue string,c *gin.Context) bool {
	//存储一周
	c.SetCookie(cname,cvalue,604800,"/","",false,true)
	return true
}

func GetCookie(cname string,c *gin.Context) string  {
	if val,err := c.Cookie(cname);err == nil{
		return val
	}
	return ""
}

func ClearCookie(cname string,c *gin.Context)  {
	c.SetCookie(cname,"",-1,"/","",false,true)
}