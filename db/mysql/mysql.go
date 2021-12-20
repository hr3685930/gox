package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlDB struct {
	database string
	host     string
	port     string
	username string
	password string
	debug    bool
}

func NewMysqlDB(database, host, port, username, password string, debug bool) *MysqlDB {
	return &MysqlDB{database, host, port, username, password, debug}
}

func (m *MysqlDB) Connect() (error, *gorm.DB) {
	dsn := m.username + ":" + m.password + "@(" + m.host + ":" + m.port + ")/" + m.database +
		"?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"

	loglevel := logger.Error
	if m.debug {
		loglevel = logger.Info
	}

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
