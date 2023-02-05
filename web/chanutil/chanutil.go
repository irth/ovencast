package chanutil

import (
	"context"
)

func Put[T any](ctx context.Context, ch chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
		return nil
	}
}

func Get[T any](ctx context.Context, ch <-chan T) (*T, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()
	case v, more := <-ch:
		return &v, more, nil
	}
}
