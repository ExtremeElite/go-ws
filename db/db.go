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
			log.Println(err)
		}
	}()
	defaultDB:= conf.CommonSet.DefaultDB
	switch defaultDB {
	default:
		DB=localMysql()
	}
}

