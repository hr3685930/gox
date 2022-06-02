package goo

import (
	"fmt"
	"github.com/pkg/errors"
)

var AsyncErr chan error
var AsyncErrFunc func(stack string)
type AsyncFunc func() error

func New() {
	AsyncErr = make(chan error, 1)
	AsyncErrFunc = errHandler
	go func() {
		for {
			select {
			case err := <-AsyncErr:
				e, ok := err.(interface{ Error })
				var stack string
				if ok {
					stack = e.GetStack()
				}else{
					stack = fmt.Sprintf("%+v\n", errors.New(err.Error()))
				}
				AsyncErrFunc(stack)
			}
		}
	}()
}

// GO 无需等待处理完成
func GO(fns AsyncFunc) {
	go func(f AsyncFunc) {
		defer func() {
			if err := recover(); err != nil {
				AsyncErr <- errors.Errorf("%+v\n", err)
			}
		}()
		err := f()
		if err != nil {
			AsyncErr <- err
		}
	}(fns)
}

func errHandler(stack string)  {
	fmt.Printf(stack)
}