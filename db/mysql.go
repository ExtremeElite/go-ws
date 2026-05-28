package db

import (
	"fmt"
	"log"
	"time"
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
		DSN:                       dsn,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}), &gormConfig)
	if err != nil {
		log.Fatal("mysql init failed: ", err)
	}
	if sqlDB, err := gormDB.DB(); err == nil {
		sqlDB.SetMaxIdleConns(localBase.MaxConnect)
		sqlDB.SetMaxOpenConns(localBase.MaxConnect * 2)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(3 * time.Minute)
	} else {
		log.Fatal("mysql build gorm failed: ", err)
	}
	return
}
