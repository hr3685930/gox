package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var ConnStore sync.Map

var Orm *gorm.DB

var DefaultLogLevel = logger.Error

func GetConnect(con string) *gorm.DB {
	v, ok := ConnStore.Load(con)
	if ok {
		return v.(*gorm.DB)
	}
	return nil
}

type DB interface {
	Connect() (error, *gorm.DB)
}
