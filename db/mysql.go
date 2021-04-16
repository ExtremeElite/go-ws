package db

import (
	"fmt"
	"log"
	"ws/common"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
)

func localMysql() *gorm.DB {
	var mysqlDB *gorm.DB
	var localBase = common.MysqlSet
	var err error
	defer func() {
		if common.Debug {
			return
		}
		if err := recover(); err != nil {
			log.Println("sql server: ", err)
			rawDB, err := mysqlDB.DB()
			if err != nil {
				log.Println("err is ", err)
				return
			}
			_ = rawDB.Close()
			return
		}
	}()
	linked := fmt.Sprintf(common.MysqlTcpConnect, localBase.User, localBase.Password, localBase.ServerHost, localBase.Port, localBase.Db)

	var loggerDefaultMode = logger.Silent
	if common.Debug {
		loggerDefaultMode=logger.Info
	}
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(loggerDefaultMode),
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
	if err != nil {
		panic(err)
	}
	sqlDbRaw, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}
	sqlDbRaw.SetMaxIdleConns(localBase.MaxConnect)
	sqlDbRaw.SetMaxOpenConns(localBase.MaxConnect * 2)
	sqlDbRaw.SetConnMaxLifetime(-1)
	return mysqlDB
}
