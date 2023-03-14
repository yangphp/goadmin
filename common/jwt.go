package common

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"strconv"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims //嵌套
	UserName   string
	UserId int
}
//jwtkey 确定后不要修改
var JWTKEY = []byte("02qRNfucOf5jYQbE")

// 获取token
func GetToken(name string,uid int) string {
	//payload
	claims := UserClaims{
		UserName: name,
		UserId:uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3600 * time.Second).Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),                        //签发时间
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(JWTKEY)

	if err != nil {

		log.Fatal(err.Error())
	}
	//将token加入到redis
	idstr := strconv.Itoa(uid)
	SetRedisVal("token"+idstr,tokenString,3500)

	return tokenString
}