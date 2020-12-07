package db

import (
	"fmt"
	"log"
	"os"
	"time"
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
	defaultLog:=logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		Colorful:      false,
		LogLevel:      logger.Error,
	})
	gormConfig := gorm.Config{
		Logger:defaultLog ,
	}
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
