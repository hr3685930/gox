package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type SqliteDB struct {
	debug    bool
}

func NewSqliteDB(debug bool) *SqliteDB {
	return &SqliteDB{debug}
}

func (m *SqliteDB) Connect() (error, *gorm.DB) {
	loglevel := logger.Error
	if m.debug {
		loglevel = logger.Info
	}
	orm, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
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
