package broadcast

import (
	"context"
)

func put[T any](ctx context.Context, ch chan T, v T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- v:
		return nil
	}
}

func get[T any](ctx context.Context, ch chan T) (*T, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case v := <-ch:
		return &v, nil
	}
}
