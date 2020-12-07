package db

import (
	"gorm.io/gorm"
	"log"
	"ws/conf"
)

var DB *gorm.DB

func init() {
	defer func() {
		if err:=recover();err!=nil {
			log.Fatal("sql server: ",err)
			return
		}
		rawDB,err:=DB.DB()
		if err!=nil {
			log.Println("err is ",err)
			return
		}
		_ = rawDB.Close()
	}()
	defaultDB:= conf.CommonSet.DefaultDB
	switch defaultDB {
	default:
		DB=localMysql()
	}
}

