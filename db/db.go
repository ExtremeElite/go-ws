package db

import (
	"gorm.io/gorm"
	"ws/conf"
)

var DB *gorm.DB

func init() {
	defaultDB:= conf.CommonSet.DefaultDB
	switch defaultDB {
	default:
		DB=localMysql()
	}
}

