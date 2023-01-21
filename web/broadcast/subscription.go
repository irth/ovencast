package broadcast

type Subscription[T any] struct {
	out   chan T
	id    int64
	unsub chan unsub[T]
}

func (s Subscription[T]) Ch() <-chan T {
	return s.out
}

func (s Subscription[T]) Unsubscribe() {
	sync := make(chan bool, 1)
	s.unsub <- unsub[T]{
		id:   s.id,
		sync: sync,
	}
	<-sync
}
