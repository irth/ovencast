package chanutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/irth/ovencast/web/chanutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	ch := make(chan int, 1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err := chanutil.Get(ctx, ch)
	assert.ErrorIs(t, err, context.Canceled)

	ctx, cancel = context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()
	<-ctx.Done()

	_, _, err = chanutil.Get(ctx, ch)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch <- 5
	msg, more, err := chanutil.Get(ctx, ch)
	require.NoError(t, err)
	require.True(t, more)
	require.NotNil(t, msg)
	require.Equal(t, 5, *msg)

	go func() {
		<-time.After(200 * time.Millisecond)
		close(ch)
	}()

	// channel close happening while Get is running
	_, more, err = chanutil.Get(ctx, ch)
	require.NoError(t, err)
	require.False(t, more)

	// channel close before Get is called
	_, more, err = chanutil.Get(ctx, ch)
	require.NoError(t, err)
	require.False(t, more)
}

func TestPut(t *testing.T) {
	ch := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := chanutil.Put(ctx, ch, 5)
	assert.ErrorIs(t, err, context.Canceled)

	ctx, cancel = context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()
	<-ctx.Done()

	err = chanutil.Put(ctx, ch, 5)
	assert.ErrorIs(t, err, context.DeadlineExceeded)

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chNonBlocking := make(chan int, 1)
	err = chanutil.Put(ctx, chNonBlocking, 5)
	require.NoError(t, err)
	require.Equal(t, 5, <-chNonBlocking)
}
