package config

import (
	"reflect"
	"sts/pkg/db"
	"sts/pkg/db/mysql"
)

type MYSQLDrive struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	App      App
}

func (m MYSQLDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		switch typeInfo.Field(i).Name {
		case "Name":
			m.Host = valInfo.Field(i).String()
			break
		case "Port":
			m.Port = valInfo.Field(i).String()
			break
		case "Database":
			m.Database = valInfo.Field(i).String()
			break
		case "Username":
			m.Username = valInfo.Field(i).String()
			break
		case "Password":
			m.Password = valInfo.Field(i).String()
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

	mysqlDB := mysql.NewMysqlDB(m.Database, m.Host, m.Port, m.Username, m.Password, m.App.Debug)
	err, orm := mysqlDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}


func (m MYSQLDrive) Default(key string) {
	db.Orm = db.GetConnect(key)
}