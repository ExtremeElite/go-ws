package db

import (
	"gorm.io/gorm"
	"ws/common"
)

var DB *gorm.DB

func init() {
	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Println("sql server: ", err)
	//		rawDB, err := DB.DB()
	//		if err != nil {
	//			log.Println("err is ", err)
	//			return
	//		}
	//		_ = rawDB.Close()
	//		return
	//	}
	//}()
	defaultDB := common.Setting.DefaultDB
	switch defaultDB {
	default:
		DB = localMysql()
	}
}
