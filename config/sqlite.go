package config

import (
	"github.com/hr3685930/pkg/db"
	"github.com/hr3685930/pkg/db/sqlite"
	"reflect"
)

type SQLiteDrive struct {
	App      App
}

func (s SQLiteDrive) Connect(key string, options interface{}, app interface{}) error {
	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			s.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			s.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			s.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}

	sqliteDB := sqlite.NewSqliteDB(s.App.Debug)
	err, orm := sqliteDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}


func (m SQLiteDrive) Default(key string) {
	db.Orm = db.GetConnect(key)
}