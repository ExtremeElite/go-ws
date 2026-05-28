package db

import (
	"ws/common"

	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	defaultDB := common.DB.Default
	switch defaultDB {
	case "mysql":
		GormDB = localMysql()
	default:
		GormDB = localMysql()
	}
}
