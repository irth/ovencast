package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type StreamState struct {
	Online bool `json:"online"`
}

func (s StreamState) Equal(other StreamState) bool {
	return s.Online == other.Online
}

// fetchState calls OME APIs to determine the stream state and returns it as a
// StreamState value. It does not modify the saved state, see
// updateStateAndNotify for that.
func (a *API) fetchState() (StreamState, error) {
	online, err := a.OME.StreamExists("default", "live", "stream")
	if err != nil {
		return StreamState{}, err
	}

	return StreamState{
		Online: online,
	}, nil
}

// updateStateAndNotify obtains a new StreamState value and saves it in the API
// struct. If it has changed, it broadcasts the update to the clients.  Returns
// true if the state has changed. If err != nil, the return value is always
// false.
func (a *API) updateStateAndNotify() (changed bool, err error) {
	state, err := a.fetchState()
	if err != nil {
		return false, fmt.Errorf("couldn't update stream state: %w", err)
	}

	changed = !a.state.Equal(state)
	if changed {
		err = a.stateUpdates.Broadcast(context.TODO(), state)
		if err != nil {
			return false, fmt.Errorf("couldn't broadcast state change: %w", err)
		}
		a.state = state
	}

	return changed, nil
}

// runStateUpdater starts a loop that keeps track of the stream state changes
// regularly. The state is checked at a slower interval, until we get notified
// by OME that a stream is starting/ending. Then, the checks happen more often,
// until the state changes or we hit a timeout.
func (a *API) runStateUpdater() {
	a.updateStateAndNotify()

	changing := false
	deadline := time.Now()

	regularInterval := 5 * time.Second
	changingInterval := 200 * time.Millisecond
	changingTimeout := 5 * time.Second

	interval := regularInterval

	for {
		log.Printf("checking state (starting: %v, interval: %v)", changing, interval)

		select {
		case <-time.After(interval):
		case <-a.admissionWebhookSignal:
			// we got notified, let's check more often
			changing = true
			interval = changingInterval
			deadline = time.Now().Add(changingTimeout)
		}

		// try to update the state & notify clients
		changed, err := a.updateStateAndNotify()
		if err != nil {
			log.Println("while updating state: %w", err)
			continue
		}

		if changing && (changed || time.Now().After(deadline)) {
			// the change got through to us, or the timeout has been hit, switch
			// back to the regular interval
			changed = false
			interval = regularInterval
		}
	}
}
