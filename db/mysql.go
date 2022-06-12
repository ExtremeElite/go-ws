package db

import (
	"fmt"
	"log"
	"ws/common"
	"ws/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func localMysql() (gormDB *gorm.DB) {
	var localBase = common.MysqlSet
	var err error
	linked := fmt.Sprintf(util.MysqlTcpConnect, localBase.User, localBase.Password, localBase.ServerHost, localBase.Port, localBase.Db)

	var loggerDefaultMode = logger.Silent
	if common.Debug {
		loggerDefaultMode = logger.Info
	}
	gormConfig := gorm.Config{
		Logger: logger.Default.LogMode(loggerDefaultMode),
	}
	gormDB, err = gorm.Open(mysql.New(mysql.Config{
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
	if sqlDB, err := gormDB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(localBase.MaxConnect)
		sqlDB.SetMaxOpenConns(localBase.MaxConnect * 2)
		sqlDB.SetConnMaxLifetime(-1)
	} else {
		log.Println("mysql build gorm failed:", err.Error())
	}
	return
}
