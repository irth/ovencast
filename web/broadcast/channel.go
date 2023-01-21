package broadcast

import (
	"context"
	"log"
	"sync"
)

type sub[T any] struct {
	out  chan T
	id   chan int64
	sync chan bool
}

type unsub[T any] struct {
	id   int64
	sync chan bool
}

type subcount chan int

type broadcast[T any] struct {
	msg T
	ctx context.Context
}

type Channel[T any] struct {
	in chan broadcast[T]

	sub      chan sub[T]
	unsub    chan unsub[T]
	subcount chan subcount

	subs    map[int64]chan T
	counter int64

	closeOnce sync.Once
}

func NewChannel[T any](ctx context.Context) *Channel[T] {
	ch := &Channel[T]{
		in: make(chan broadcast[T], 8),

		sub:      make(chan sub[T], 8),
		unsub:    make(chan unsub[T], 8),
		subcount: make(chan subcount, 8),

		subs:    make(map[int64]chan T, 64),
		counter: 0,
	}
	go ch.handle(ctx)
	return ch
}

func (c *Channel[T]) handle(ctx context.Context) {
	defer func() {
		for _, sub := range c.subs {
			close(sub)
		}
	}()

	for {
		log.Printf("select")
		select {
		case <-ctx.Done():
			log.Println("im outta here")
			return
		case sub := <-c.sub:
			id := c.subscribe(sub.out)
			sub.id <- id
			sub.sync <- true
		case unsub := <-c.unsub:
			c.unsubscribe(unsub.id)
			unsub.sync <- true
		case subcount := <-c.subcount:
			log.Printf("subcount")
			subcount <- len(c.subs)
			log.Printf("subcount done")
		case broadcast, more := <-c.in:
			if !more {
				return
			}
			c.broadcast(broadcast)
		}
	}
}

func (c *Channel[T]) subscribe(out chan T) int64 {
	id := c.counter
	log.Printf("sub: %d", id)
	c.counter += 1
	c.subs[id] = out
	return id
}

func (c *Channel[T]) unsubscribe(id int64) {
	ch, ok := c.subs[id]
	if !ok {
		return
	}
	delete(c.subs, id)
	close(ch)
}

func (c *Channel[T]) broadcast(b broadcast[T]) {
	resend := make(map[int64]chan T, len(c.subs))

	// try to send without blocking (will work if queues are empty)
	for id, sub := range c.subs {
		log.Printf("broadcast: sending (1, %d)", id)
		select {
		case <-b.ctx.Done():
			log.Printf("broadcast: ctx expired (1)")
			return // oops, context expired
		case sub <- b.msg:
			log.Printf("broadcast: sent (1, %d)", id)
			continue // sent successfuly
		default:
			// we failed.
			log.Printf("broadcast: blocked (1, %d)", id)
			resend[id] = sub
		}
	}

	// now send in a blocking manner (a blocked send might cause the next
	// subscribers to also not receive the message if the context expires)

	success := make(map[int64]bool, len(c.subs))
	ctxExpired := false
outer:
	for id, sub := range resend {
		log.Printf("broadcast: sending (2, %d)", id)
		select {
		case <-b.ctx.Done():
			log.Println("broadcast: ctx expired (2)")
			ctxExpired = true
			break outer
		case sub <- b.msg:
			log.Printf("broadcast: sent (2, %d)", id)
			success[id] = true
			continue
		}
	}

	// unsubscribe missbehaving
	if ctxExpired {
		for id, sub := range resend {
			log.Printf("broadcast: checking (%d)", id)
			select {
			case sub <- b.msg:
				log.Printf("broadcast: sent (2, %d)", id)
				success[id] = true
				continue
			default:
				log.Printf("unsubscribe: %d", id)
				c.unsubscribe(id)
			}
		}
	}
}

func (c *Channel[T]) Broadcast(ctx context.Context, message T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.in <- broadcast[T]{
		ctx: ctx,
		msg: message,
	}:
		return nil
	}
}

func (c *Channel[T]) Close() {
	c.closeOnce.Do(func() { close(c.in) })
}

func (c *Channel[T]) Subscribe() Subscription[T] {
	// TODO: add context support here
	s := Subscription[T]{out: make(chan T, 16)}

	sync := make(chan bool, 1)
	id := make(chan int64, 1)

	c.sub <- sub[T]{
		out:  s.out,
		sync: sync,
		id:   id,
	}

	<-sync
	s.id = <-id

	return s
}

func (c *Channel[T]) SubscriberCount(ctx context.Context) (int, error) {
	subcount := make(subcount, 1)
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case c.subcount <- subcount:
		// request sent
	}
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case s := <-subcount:
		return s, nil
	}
}
