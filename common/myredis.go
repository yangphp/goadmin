package common

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

var Redi redis.Conn
var err error
func init()  {
	//连接redis
	Redi, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
}

func SetRedisVal(name string,value string,expire int) bool {
	Redi.Do("set",name,value)
	Redi.Do("expire",name,expire)
	return true
}

func GetRedisVal(name string) string  {

	res,err := Redi.Do("get",name)
	if err != nil || res == nil {
		return ""
	}

	val,err := redis.String(res,err)
	if err != nil {
		return ""
	}

	return val
}





