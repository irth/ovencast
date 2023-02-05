package broadcast_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/irth/ovencast/web/broadcast"
	"github.com/irth/ovencast/web/chanutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func TestSubscribe(t *testing.T) {
	ctx, cancel := testctx()
	defer cancel()

	ch := broadcast.NewChannel[int]()
	go ch.Run(ctx)

	subcount, err := ch.SubCount(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, subcount)

	_, err = ch.Subscribe(ctx)
	assert.NoError(t, err)

	subcount, err = ch.SubCount(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, subcount)

	_, err = ch.Subscribe(ctx)
	assert.NoError(t, err)

	_, err = ch.Subscribe(ctx)
	assert.NoError(t, err)

	subcount, err = ch.SubCount(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, subcount)
}

func startSubscribers[T any](t testing.TB, ctx context.Context, ch *broadcast.Channel[T], n int) chan T {
	returnCh := make(chan T)

	for i := 0; i < n; i += 1 {
		sub, err := ch.Subscribe(ctx)
		assert.NoError(t, err)

		s := i
		go func() {
			getNo := 0
			for {
				m, more, err := chanutil.Get(ctx, sub)
				require.NoError(t, err, "get from subscriber %d (%d)", s, getNo)
				getNo += 1
				if !more {
					return
				}
				require.NotNil(t, m)
				chanutil.Put(ctx, returnCh, *m)
			}
		}()
	}

	subCount, err := ch.SubCount(ctx)
	require.NoError(t, err)
	require.Equal(t, n, subCount)

	return returnCh
}

func TestBroadcast(t *testing.T) {
	ctx, cancel := testctx()
	defer cancel()

	ch := broadcast.NewChannel[int]()
	defer ch.Close(ctx)
	go ch.Run(ctx)

	returned := 0
	n := 8

	returnCh := startSubscribers(t, ctx, ch, n)

	go func() {
		err := chanutil.Put(ctx, ch.Ch(), 42)
		require.NoError(t, err, "put received an error")
	}()

	for returned < n {
		ret, _, err := chanutil.Get(ctx, returnCh)
		require.NoError(t, err)
		require.Equal(t, 42, *ret)
		returned += 1
		fmt.Println("get", returned)
	}

	require.Equal(t, n, returned)
}

func BenchmarkBroadcast(b *testing.B) {
	ctx, cancel := testctx()
	defer cancel()

	ch := broadcast.NewChannel[int]()
	go ch.Run(ctx)
	defer ch.Close(ctx)

	subCount := 10_000

	returnCh := startSubscribers(b, ctx, ch, subCount)
	for i := 0; i < b.N; i += 1 {
		returned := 0

		err := chanutil.Put(ctx, ch.Ch(), 42)
		require.NoError(b, err, "put received an error")

		for returned < subCount {
			ret, _, err := chanutil.Get(ctx, returnCh)
			require.NoError(b, err)
			require.Equal(b, 42, *ret)
			returned += 1
		}
	}
}

func TestClose(t *testing.T) {
	ctx, cancel := testctx()
	defer cancel()

	ch := broadcast.NewChannel[int]()
	go ch.Run(ctx)

	sub1, err := ch.Subscribe(ctx)
	assert.NoError(t, err)

	sub2, err := ch.Subscribe(ctx)
	assert.NoError(t, err)

	err = ch.Close(ctx)
	assert.NoError(t, err)

	_, more, err := chanutil.Get(ctx, sub1)
	assert.NoError(t, err)
	assert.False(t, more)

	_, more, err = chanutil.Get(ctx, sub2)
	assert.NoError(t, err)
	assert.False(t, more)
}
