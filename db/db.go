package db

import (
	"ws/common"

	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	defaultDB := common.Setting.DefaultDB
	switch defaultDB {
	default:
		GormDB = localMysql()
	}
}
