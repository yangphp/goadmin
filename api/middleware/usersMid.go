package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goadmin/common"
	"net/http"
	"strconv"
)

func CheckAuth(c *gin.Context)  {

	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "no token"})
		c.Abort() //终止
		return
	}

	claims := common.UserClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return common.JWTKEY, nil
	})
	//验证失败
	if err != nil {
		ve, _ := err.(*jwt.ValidationError)
		if ve.Errors == jwt.ValidationErrorExpired {
			c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "token expired"})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 405, "msg": "token invalid"})
		}
		c.Abort() //终止
		return
	}

	//获取redistoken，进行验证
	idstr := strconv.Itoa(claims.UserId)
	red_token := common.GetRedisVal("token"+idstr)

	if red_token != token {
		c.JSON(http.StatusOK, gin.H{"code": 406, "msg": "token not exist"})
		c.Abort() //终止
		return
	}


	c.Set("cliaims", claims) //传递参数
}

