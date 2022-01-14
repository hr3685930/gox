package goo

import (
	"fmt"
	"github.com/pkg/errors"
)

var AsyncErr chan error
var AsyncErrFunc func(err error)
type AsyncFunc func() error

func New() {
	AsyncErr = make(chan error, 1)
	AsyncErrFunc = errHandler
	go func() {
		for {
			select {
			case err := <-AsyncErr:
				AsyncErrFunc(err)
			}
		}
	}()
}

func GO(fns AsyncFunc) {
	go func(f AsyncFunc) {
		defer func() {
			if err := recover(); err != nil {
				AsyncErr <- errors.Errorf("%+v\n", err)
			}
		}()
		err := f()
		AsyncErr <- err
	}(fns)
}

func errHandler(err error)  {
	fmt.Printf("%+v\n", err)
}