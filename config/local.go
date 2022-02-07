package config

import (
	"github.com/hr3685930/pkg/queue"
	"github.com/hr3685930/pkg/queue/local"
	"reflect"
)

type LocalDrive struct {
	App App
}

func (m LocalDrive) Connect(key string, options interface{}, app interface{}) error {
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
	c := local.NewLocalMQ()
	queue.QueueStore.Store(key, c)
	return nil
}

func (r LocalDrive) Default(key string) {
	queue.MQ = queue.GetQueueDrive(key)
}
