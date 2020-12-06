package db

import (
	"fmt"
	"log"
	"ws/conf"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
)

func localMysql() *gorm.DB {
	var mysqlDB *gorm.DB
	var bs = conf.Config()
	var localBase = bs.MysqlDB
	var err error
	linked := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local`, localBase.User, localBase.Password, localBase.ServerHost, localBase.Port, localBase.Db)
	gormConfig := gorm.Config{}
	if bs.Common.Env == "dev" {
		gormConfig = gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	}
	mysqlDB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName:                "",
		DSN:                       linked,
		Conn:                      nil,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         255,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
	}), &gormConfig)
	if err!=nil {
		log.Panicln("mysql err is:", err.Error())
	}
	sqlDbRaw, err := mysqlDB.DB()
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	sqlDbRaw.SetMaxIdleConns(localBase.MaxConnect)
	sqlDbRaw.SetMaxOpenConns(localBase.MaxConnect * 2)
	sqlDbRaw.SetConnMaxLifetime(-1)
	return mysqlDB
}
