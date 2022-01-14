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

type response struct {
	res interface{}
	err error
}

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
			defer wg.Done()

			var r response
			r.res, r.err = f(ctx)

			rs[index] = r.res
			errs[index] = r.err
		}(i, fn)
	}

	wg.Wait()

	return rs, errs
}
