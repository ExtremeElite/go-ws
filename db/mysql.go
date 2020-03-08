package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"ws/conf"

	_ "github.com/go-sql-driver/mysql"
)

func mysql() *gorm.DB {
	var mysqlDB *gorm.DB
	var bs= conf.GetConfig()
	var localBase=bs.MysqlDB
	var err error
	linked:=fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`,localBase.User,localBase.Password,localBase.ServerHost,localBase.Port,localBase.Db)
	mysqlDB, err = gorm.Open("mysql", linked)
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	mysqlDB.SingularTable(true)
	mysqlDB.DB().SetMaxIdleConns(localBase.MaxConnect)
	mysqlDB.DB().SetMaxOpenConns(localBase.MaxConnect*2)
	return mysqlDB
}
