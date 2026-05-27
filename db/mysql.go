package db

import (
	"fmt"
	"log"
	"ws/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func localMysql() (gormDB *gorm.DB) {
	var localBase = common.DB.Mysql
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		localBase.User, localBase.Password, localBase.ServerHost, localBase.Port, localBase.Db)

	var loggerDefaultMode = logger.Silent
	if common.Debug {
		loggerDefaultMode = logger.Info
	}
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(loggerDefaultMode),
	}
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
		DriverName:                "",
		DSN:                       dsn,
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
	if sqlDB, err := gormDB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(localBase.MaxConnect)
		sqlDB.SetMaxOpenConns(localBase.MaxConnect * 2)
		sqlDB.SetConnMaxLifetime(-1)
	} else {
		log.Println("mysql build gorm failed:", err.Error())
	}
	return
}
