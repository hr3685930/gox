package goo

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

// Group Group
type Group struct {
	wg         sync.WaitGroup
	ResultChan chan interface{}
	ErrorChan  chan error
	Result     []interface{}
	Error      []error
}

// NewGroup num控制协程处理数
func NewGroup(num int) *Group {
	g := &Group{
		ResultChan: make(chan interface{}, num),
		ErrorChan:  make(chan error, num),
		Result:     make([]interface{}, 0),
		Error:      make([]error, 0),
	}
	go func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for err := range g.ErrorChan {
				g.Error = append(g.Error, err)
			}
		}()

		go func() {
			defer wg.Done()
			for r := range g.ResultChan {
				g.Result = append(g.Result, r)
			}
		}()
		wg.Wait()
	}()
	return g
}

// One One
func (g *Group) One(ctx context.Context, fn SyncFunc) {
	g.wg.Add(1)
	go func(f SyncFunc) {
		defer func() {
			if err := recover(); err != nil {
				g.ResultChan <- nil
				g.ErrorChan <- errors.Errorf("%+v\n", err)
			}
			g.wg.Done()
		}()
		res, err := f(ctx)
		g.ResultChan <- res
		g.ErrorChan <- err
	}(fn)
}

// Wait 返回结果为无序
func (g *Group) Wait() ([]interface{}, []error) {
	g.wg.Wait()
	close(g.ResultChan)
	close(g.ErrorChan)
	return g.Result, g.Error
}

// SyncFunc SyncFunc
type SyncFunc func(ctx context.Context) (interface{}, error)

// All 有序返回结果 func协程一次处理  error nil也返回
func All(ctx context.Context, fns ...SyncFunc) ([]interface{}, []error) {
	rs := make([]interface{}, len(fns))
	errs := make([]error, len(fns))

	var wg sync.WaitGroup
	wg.Add(len(fns))

	for i, fn := range fns {
		go func(index int, f SyncFunc) {
			defer func() {
				if err := recover(); err != nil {
					rs[index] = nil
					errs[index] = errors.Errorf("%+v\n", err)
				}
				wg.Done()
			}()

			res, err := f(ctx)
			rs[index] = res
			errs[index] = err
		}(i, fn)
	}

	wg.Wait()

	return rs, errs
}
