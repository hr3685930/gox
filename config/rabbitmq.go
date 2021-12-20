package config

import (
	"reflect"
	"github.com/hr3685930/pkg/queue"
	"github.com/hr3685930/pkg/queue/rabbitmq"
)

type RabbitMQDrive struct {
	Host     string
	Port     string
	VHost    string
	Username string
	Password string
	App      App
}

func (m RabbitMQDrive) Connect(key string, options interface{}, app interface{}) error {
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
		case "VHost":
			m.VHost = valInfo.Field(i).String()
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
	rabbitMQ := rabbitmq.NewRabbitMQ(m.Username, m.Password, m.Host, m.Port, m.VHost, m.App.Name)
	err := rabbitMQ.Connect()
	if err != nil {
		return err
	}
	queue.QueueStore.Store(key, rabbitMQ)
	return nil
}


func (r RabbitMQDrive) Default(key string) {
	queue.MQ = queue.GetQueueDrive(key)
}
