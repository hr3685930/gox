package mysql

import (
	"github.com/hr3685930/pkg/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlDB struct {
	dsn   string
	debug bool
}

func NewMysqlDB(dsn string, debug bool) *MysqlDB {
	return &MysqlDB{dsn, debug}
}

func (m *MysqlDB) Connect() (error, *gorm.DB) {
	dsn := m.dsn
	loglevel := db.DefaultLogLevel
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(loglevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return err, nil
	}
	sqlDB, _ := orm.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil, orm
}
