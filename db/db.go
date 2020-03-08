package db

import (
	"github.com/jinzhu/gorm"
	"ws/conf"
)

var DB *gorm.DB

func init() {
	defaultDB:= conf.GetConfig().DefaultDB
	switch defaultDB {
	default:
		DB=mysql()
	}
}

