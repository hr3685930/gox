package config

import (
	"github.com/hr3685930/pkg/db"
	"github.com/hr3685930/pkg/db/clickhouse"
	"reflect"
)

type ClickhouseDrive struct {
	Dsn     string
	App      App
}

func (m ClickhouseDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		switch typeInfo.Field(i).Name {
		case "Dsn":
			m.Dsn = valInfo.Field(i).String()
			break
		}
	}

	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			m.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			m.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			m.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}

	clickhouseDB := clickhouse.NewClickHouseDB(m.Dsn, m.App.Debug)
	err, orm := clickhouseDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}


func (m ClickhouseDrive) Default(key string) {
	db.Orm = db.GetConnect(key)
}