package broadcast

import (
	"context"
	"fmt"
	"sync"

	"github.com/irth/ovencast/web/chanutil"
)

type _nothing struct{}

var nothing _nothing = struct{}{}

type Channel[T any] struct {
	broadcastCh chan T

	subReq      chan *request[chan T, _nothing]
	subCountReq chan *request[_nothing, int]
	closeReq    chan *request[_nothing, _nothing]

	subs map[chan T]struct{}

	running  bool
	runningL sync.Mutex
}

func NewChannel[T any]() *Channel[T] {
	return &Channel[T]{
		// TODO: make queue sizes configurable?
		broadcastCh: make(chan T, 128),

		subReq:      make(chan *request[chan T, _nothing], 8),
		subCountReq: make(chan *request[_nothing, int], 8),
		closeReq:    make(chan *request[_nothing, _nothing], 1),

		subs: make(map[chan T]struct{}, 8),

		running: false,
	}
}

func (c *Channel[T]) Run(ctx context.Context) {
	c.runningL.Lock()
	c.running = true
	c.runningL.Unlock()

	defer func() {
		c.runningL.Lock()
		defer c.runningL.Unlock()
		c.running = false
	}()

	var closeReq *request[_nothing, _nothing] = nil
	defer func() {
		for ch := range c.subs {
			close(ch)
		}
		if closeReq != nil {
			closeReq.Ok(nothing)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case r := <-c.subReq:
			ch := r.Args()
			c.subs[ch] = nothing
			r.Ok(nothing)
			// TODO: do something with errors from r.Ok
		case r := <-c.subCountReq:
			r.Ok(len(c.subs))
		case r := <-c.closeReq:
			closeReq = r
			return
		case m := <-c.broadcastCh:
			c.broadcast(m)
		}
	}
}

func (c *Channel[T]) Ch() chan<- T {
	return c.broadcastCh
}

func (c *Channel[T]) Broadcast(ctx context.Context, m T) error {
	return chanutil.Put(ctx, c.broadcastCh, m)
}

func (c *Channel[T]) broadcast(m T) {
	for sub := range c.subs {
		select {
		case sub <- m:
		default:
			close(sub)
			delete(c.subs, sub)
		}
	}
}

func (c *Channel[T]) Subscribe(ctx context.Context) (<-chan T, error) {
	ch := make(chan T, 8)
	_, err := sendRequest(ctx, c.subReq, ch)
	if err != nil {
		return nil, fmt.Errorf("subscription failed: %w", err)
	}

	return ch, nil
}

func (c *Channel[T]) SubCount(ctx context.Context) (int, error) {
	subCount, err := sendRequest(ctx, c.subCountReq, nothing)

	if err != nil {
		return 0, fmt.Errorf("subscription failed: %w", err)
	}

	return *subCount, nil
}

func (c *Channel[T]) Close(ctx context.Context) error {
	c.runningL.Lock()
	defer c.runningL.Unlock()

	if c.running {
		_, err := sendRequest(ctx, c.closeReq, nothing)
		return err
	}

	return nil
}

func (c *Channel[T]) Running() bool {
	c.runningL.Lock()
	defer c.runningL.Unlock()

	return c.running
}
