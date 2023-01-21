package broadcast_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/irth/ovencast/web/broadcast"
	"github.com/stretchr/testify/require"
)

func TestBroadcast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ch := broadcast.NewChannel[int](ctx)
	defer ch.Close()

	s, err := ch.SubscriberCount(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, s)

	s1 := ch.Subscribe()
	s2 := ch.Subscribe()

	s, err = ch.SubscriberCount(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, s)

	ch.Broadcast(ctx, 2137)
	require.Equal(t, 2137, <-s1.Ch())
	require.Equal(t, 2137, <-s2.Ch())
}

func TestMisbehavingListenersGetUnsubscribed(t *testing.T) {
	testCtx, cancelTest := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelTest()

	ch := broadcast.NewChannel[int](testCtx)
	defer ch.Close()

	s, err := ch.SubscriberCount(testCtx)
	require.NoError(t, err)
	require.Equal(t, 0, s)

	// subscribers that don't read
	ch.Subscribe()
	ch.Subscribe()

	// a well-behaved subscriber
	sub := ch.Subscribe()
	go func() {
		log.Println("started consuming")
		for m := range sub.Ch() {
			log.Println("message:", m)
		}
		log.Println("channel closed")
	}()

	s, err = ch.SubscriberCount(testCtx)
	require.NoError(t, err)
	require.Equal(t, 3, s)

	broadcastCtx, cancelBroadcast := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelBroadcast()
	err = nil
	for err == nil {
		// fill queues
		err = ch.Broadcast(broadcastCtx, 2137)
	}
	require.ErrorIs(t, broadcastCtx.Err(), context.DeadlineExceeded)

	s, err = ch.SubscriberCount(testCtx)
	require.NoError(t, err)
	require.Equal(t, 1, s)
}
