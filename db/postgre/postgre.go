package postgre

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type PostgreDB struct {
	dsn   string
	debug bool
}

func NewPostgreDB(dsn string, debug bool) *PostgreDB {
	return &PostgreDB{dsn, debug}
}

func (m *PostgreDB) Connect() (error, *gorm.DB) {
	dsn := m.dsn
	loglevel := logger.Error
	if m.debug {
		loglevel = logger.Info
	}

	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
