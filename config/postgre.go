package config

import (
	"github.com/hr3685930/pkg/db"
	"github.com/hr3685930/pkg/db/postgre"
	"reflect"
)

type PostgreDrive struct {
	Dsn     string
	App      App
}

func (p PostgreDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		switch typeInfo.Field(i).Name {
		case "Dsn":
			p.Dsn = valInfo.Field(i).String()
			break
		}
	}

	var appTypeInfo = reflect.TypeOf(app)
	var appValInfo = reflect.ValueOf(app)
	for i := 0; i < appTypeInfo.NumField(); i++ {
		switch appTypeInfo.Field(i).Name {
		case "Name":
			p.App.Name = appValInfo.Field(i).String()
			break
		case "Env":
			p.App.Env = appValInfo.Field(i).String()
			break
		case "Debug":
			p.App.Debug = appValInfo.Field(i).Bool()
			break
		}
	}

	postgreDB := postgre.NewPostgreDB(p.Dsn, p.App.Debug)
	err, orm := postgreDB.Connect()
	if err != nil {
		return err
	}
	db.ConnStore.Store(key, orm)
	return nil
}


func (m PostgreDrive) Default(key string) {
	db.Orm = db.GetConnect(key)
}