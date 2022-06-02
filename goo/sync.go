package goo

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type Group struct {
	wg     sync.WaitGroup
	Result chan interface{}
	Error  chan error
}

// NewGroup num控制协程处理数
func NewGroup(num int) *Group {
	rs := make(chan interface{}, num)
	err := make(chan error, num)
	g := &Group{Result: rs, Error: err}
	return g
}

func (g *Group) One(ctx context.Context, fn SyncFunc) {
	g.wg.Add(1)
	go func(f SyncFunc) {
		defer func() {
			if err := recover(); err != nil {
				g.Result <- nil
				g.Error <- errors.Errorf("%+v\n", err)
			}
			g.wg.Done()
		}()
		res, err := f(ctx)
		g.Result <- res
		g.Error <- err
	}(fn)
}

// Wait 返回结果为无序
func (g *Group) Wait() ([]interface{}, []error) {
	g.wg.Wait()
	rs := make([]interface{}, 0)
	errs := make([]error, 0)
	close(g.Result)
	close(g.Error)

	var w sync.WaitGroup
	w.Add(2)
	go func() {
		defer w.Done()
		for err := range g.Error {
			errs = append(errs, err)
		}
	}()

	go func() {
		defer w.Done()
		for r := range g.Result {
			rs = append(rs, r)
		}
	}()
	w.Wait()
	return rs, errs
}

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
