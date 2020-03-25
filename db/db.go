package db

import (
	"github.com/jinzhu/gorm"
	"ws/conf"
)

var DB *gorm.DB

func init() {
	defaultDB:= conf.CommonSet.DefaultDB
	switch defaultDB {
	default:
		DB=mysql()
	}
}

