package broadcast

type Channel[T any] struct {
	subs map[chan T]struct{}
}
