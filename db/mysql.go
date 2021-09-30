package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"ws/common"
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
		loggerDefaultMode = logger.Info
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
		log.Println("mysql init failed:", err.Error())
	}
	if sqlDB, err := mysqlDB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(localBase.MaxConnect)
		sqlDB.SetMaxOpenConns(localBase.MaxConnect * 2)
		sqlDB.SetConnMaxLifetime(-1)
	} else {
		log.Println("mysql build gorm failed:", err.Error())
	}
	return mysqlDB
}
