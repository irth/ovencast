package broadcast

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/irth/ovencast/web/chanutil"
	"github.com/stretchr/testify/require"
)

func testctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*1)
}

// TODO: test contexts

func TestRequestCanOnlyBeAnsweredOnce(t *testing.T) {
	ctx, cancel := testctx()
	defer cancel()

	// to try out Ok then Ok
	reqCh := make(chan *request[int, int], 1) // queue to avoid blocking
	defer close(reqCh)

	errCh := make(chan error, 2) // to return errors from goroutine

	testVal1 := 21
	testVal2 := 37
	testErr1 := fmt.Errorf("err1")
	testErr2 := fmt.Errorf("err2")

	go func() {
		for {
			select {
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			case r, more := <-reqCh:
				if !more {
					return
				}
				arg := r.Args()
				switch arg {
				case 0:
					errCh <- r.Ok(testVal1)
					errCh <- r.Ok(testVal2)
				case 1:
					errCh <- r.Err(testErr1)
					errCh <- r.Err(testErr2)
				case 2:
					errCh <- r.Ok(testVal1)
					errCh <- r.Err(testErr1)
				case 3:
					errCh <- r.Err(testErr1)
					errCh <- r.Ok(testVal1)
				}
			}
		}
	}()

	checkVal := func(expected int, res *int, err error) {
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, expected, *res)
	}

	checkErrors := func(errCh chan error) {
		testErr, _, err := chanutil.Get(ctx, errCh)
		require.NoError(t, err)
		require.NotNil(t, testErr)
		require.Nil(t, *testErr)

		testErr, _, err = chanutil.Get(ctx, errCh)
		require.NoError(t, err)
		require.NotNil(t, testErr)
		require.ErrorIs(t, *testErr, ErrResultSentTwice)
	}

	res, err := sendRequest(ctx, reqCh, 0)
	checkVal(testVal1, res, err)
	checkErrors(errCh)

	_, err = sendRequest(ctx, reqCh, 1)
	require.ErrorIs(t, err, testErr1)
	checkErrors(errCh)

	res, err = sendRequest(ctx, reqCh, 2)
	checkVal(testVal1, res, err)
	checkErrors(errCh)

	_, err = sendRequest(ctx, reqCh, 3)
	require.ErrorIs(t, err, testErr1)
	checkErrors(errCh)
}
