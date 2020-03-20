package db

import (
	"github.com/jinzhu/gorm"
	"ws/conf"
)

var DB *gorm.DB

func init() {
	defaultDB:= conf.Config().Common.DefaultDB
	switch defaultDB {
	default:
		DB=mysql()
	}
}

