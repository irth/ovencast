package chanutil

import (
	"context"
	"fmt"
)

type ContextError struct {
	err error
}

func (c ContextError) Error() string {
	return fmt.Sprintf("context error: %s", c.err.Error())
}

func (c ContextError) Unwrap() error {
	return c.err
}

func Put[T any](ctx context.Context, ch chan<- T, v T) error {
	select {
	case <-ctx.Done():
		return ContextError{ctx.Err()}
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
