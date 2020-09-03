package db

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"ws/conf"

	"gorm.io/driver/mysql"
)

func localMysql() *gorm.DB {
	var mysqlDB *gorm.DB
	var bs= conf.Config()
	var localBase=bs.MysqlDB
	var err error
	linked:=fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`,localBase.User,localBase.Password,localBase.ServerHost,localBase.Port,localBase.Db)
	mysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName:                "",
		DSN:                       linked,
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         255,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
	}), &gorm.Config{

	})
	if sqlDbRaw, err := mysqlDB.DB();err!=nil{
		sqlDbRaw.SetMaxIdleConns(localBase.MaxConnect)
		sqlDbRaw.SetMaxOpenConns(localBase.MaxConnect*2)
		sqlDbRaw.SetConnMaxLifetime(-1)
		if bs.Common.Env=="dev" {
			mysqlDB.Logger.LogMode(logger.Info)
		}

	}
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	return mysqlDB
}
