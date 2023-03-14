package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB
var dbErr error

func init() {
	dsn := "root:root123456@tcp(127.0.0.1:3306)/godemo?charset=utf8"
	DB, dbErr = gorm.Open(mysql.Open(dsn),&gorm.Config{
		//全局配置
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,    //表名不增加复数
		},
	})
	sqlDb,_ := DB.DB()
	//设置连接超时
	sqlDb.SetConnMaxLifetime(3000) //设置最大生存时间，3000秒，比数据库的3600要小
	sqlDb.SetMaxIdleConns(100) //设置最大限制连接
	sqlDb.SetMaxOpenConns(100) //设置最大连接数

	if dbErr != nil {
		log.Fatal(dbErr.Error())
	}
}



