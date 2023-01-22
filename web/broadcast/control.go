package broadcast

import (
	"context"
	"fmt"
	"sync"
)

type result[T any] struct {
	ok  *T
	err error
}

func (r result[T]) Unpack() (*T, error) {
	return r.ok, r.err
}

type request[Args any, Ret any] struct {
	args Args

	ret  chan result[Ret]
	ctx  context.Context
	once sync.Once
}

func sendRequest[Args any, Ret any](ctx context.Context, ch chan *request[Args, Ret], args Args) (*Ret, error) {
	req := &request[Args, Ret]{
		ctx:  ctx,
		args: args,

		// channel queue size of 1 to make sure we don't stop execution
		// of the main loop
		ret: make(chan result[Ret], 1),
	}

	// put and get fail if context gets cancelled
	err := put(ctx, ch, req)
	if err != nil {
		return nil, err
	}

	ret, err := get(ctx, req.ret)
	if err != nil {
		return nil, err
	}

	return ret.Unpack()
}

func (r *request[Args, Ret]) Context() context.Context {
	return r.ctx
}

func (r *request[Args, Ret]) Args() Args {
	return r.args
}

var ErrResultSentTwice = fmt.Errorf("tried to send result twice")

func (r *request[Args, Ret]) sendResult(res result[Ret]) error {
	done := false
	r.once.Do(func() {
		r.ret <- res
		done = true
	})
	if !done {
		return ErrResultSentTwice
	}
	return nil
}

func (r *request[Args, Ret]) Ok(ret Ret) error {
	return r.sendResult(result[Ret]{ok: &ret})
}

func (r *request[Args, Ret]) Err(err error) error {
	return r.sendResult(result[Ret]{err: err})
}
