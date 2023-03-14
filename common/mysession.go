package common

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetSession(sname string,svalue string,c *gin.Context)  {
	// 初始化session对象
	session := sessions.Default(c)
	// 过期时间2h
	session.Options(sessions.Options{MaxAge: 3600 * 24,Path: "/"})

	session.Set(sname,svalue)
	session.Save()
}

func GetSession(sname string,c *gin.Context) string  {
	session := sessions.Default(c)


	if val,err := session.Get(sname).(string);err == true{
		return val
	}else{
		return ""
	}
}

func ClearSession(sname string,c *gin.Context)  {
	session := sessions.Default(c)
	session.Delete(sname)
	session.Clear()
	session.Save()
	return
}

