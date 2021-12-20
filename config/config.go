package config

import (
	"github.com/creasty/defaults"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type App struct {
	Name  string
	Env   string
	Debug bool
}

type Conf interface {
	Connect(key string, options interface{}, app interface{}) error
	Default(key string)
}

func Drive(driveEnv, app interface{}) error {
	var typeInfo = reflect.TypeOf(driveEnv)
	var valInfo = reflect.ValueOf(driveEnv)
	num := typeInfo.NumField()
	for i := 0; i < num; i++ {
		params := make([]reflect.Value, 3)
		params[0] = reflect.ValueOf(typeInfo.Field(i).Name)
		params[1] = valInfo.Field(i)
		params[2] = reflect.ValueOf(app)
		item := valInfo.Field(i).MethodByName("Connect")
		res := item.Call(params)
		err := res[0].Interface()
		if err != nil {
			return res[0].Interface().(error)
		}
	}
	dOption := make([]reflect.Value, 1)
	dOption[0] = reflect.ValueOf(typeInfo.Field(0).Name)
	d := valInfo.Field(0).MethodByName("Default")
	d.Call(dOption)
	return nil
}


// Load priority env > yaml > default
func Load(e interface{}) error {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transform config struct to map
	confMap := structs.Map(e)
	if err := mapstructure.Decode(e, &confMap); err != nil {
		return errors.Wrap(err, "Unable to Decode config")
	}

	// Flatten nested conf map
	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return errors.Wrap(err, "Unable to flatten config")
	}

	// Bind each conf fields to environment vars
	for key, _ := range flat {
		err := v.BindEnv(key)
		if err != nil {
			return errors.Wrapf(err, "Unable to bind env var: %s", key)
		}
	}

	_ = v.ReadInConfig()
	if err := defaults.Set(e); err != nil {
		return err
	}
	if err := v.Unmarshal(e); err != nil {
		return err
	}
	return nil
}