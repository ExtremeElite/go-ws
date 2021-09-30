package db

import (
	"gorm.io/gorm"
	"ws/common"
)

var GormDB *gorm.DB

func init() {
	defaultDB := common.Setting.DefaultDB
	switch defaultDB {
	default:
		GormDB = localMysql()
	}
}
