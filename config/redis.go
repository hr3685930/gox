package config

import (
	"reflect"
	"strconv"
	"sts/pkg/cache"
	"sts/pkg/cache/redis"
)

type RedisDrive struct {
	Host     string
	Port     string
	Database string
	Auth     string
	App      App
}

func (m RedisDrive) Connect(key string, options interface{}, app interface{}) error {
	var typeInfo = reflect.TypeOf(options)
	var valInfo = reflect.ValueOf(options)
	for i := 0; i < typeInfo.NumField(); i++ {
		switch typeInfo.Field(i).Name {
		case "Host":
			m.Host = valInfo.Field(i).String()
			break
		case "Port":
			m.Port = valInfo.Field(i).String()
			break
		case "Database":
			m.Database = valInfo.Field(i).String()
			break
		case "Auth":
			m.Auth = valInfo.Field(i).String()
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
	database, err := strconv.Atoi(m.Database)
	if err != nil {
		return err
	}
	c, err := redis.New(m.Host, m.Port, database, m.Auth)
	if err != nil {
		return err
	}
	cache.CacheMap.Store(key, c)
	return nil
}


func (r RedisDrive) Default(key string) {
	cache.Cached = cache.GetCache(key)
}