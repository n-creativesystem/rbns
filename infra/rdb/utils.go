package rdb

import (
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"gorm.io/gorm"
)

func GetLogger(db *gorm.DB) logger.Logger {
	var log logger.Logger
	if gormLog, ok := db.Logger.(*gormLogger); ok {
		log = gormLog.log
	} else {
		log = logger.GetLog()
	}
	return log
}
