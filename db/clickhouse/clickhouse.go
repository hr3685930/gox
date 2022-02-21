package clickhouse

import (
	"github.com/hr3685930/pkg/db"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type ClickHouseDB struct {
	dsn   string
	debug bool
}

func NewClickHouseDB(dsn string, debug bool) *ClickHouseDB {
	return &ClickHouseDB{dsn, debug}
}

func (c *ClickHouseDB) Connect() (error, *gorm.DB) {
	dsn := c.dsn
	loglevel := db.DefaultLogLevel
	if c.debug {
		loglevel = logger.Info
	}

	orm, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
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
