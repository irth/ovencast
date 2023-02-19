package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/irth/wsrpc"
)

// State (GET <api>/state) fetches the up-to-date stream state from OME and
// returns it to the client.
func (a *API) State(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("ws") == "1" {
		a.StateWS(w, r)
		return
	}

	state, err := a.fetchState()

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response[Empty]{
			OK:    false,
			Error: err.Error(),
		})
	}

	json.NewEncoder(w).Encode(Response[StreamState]{
		OK:       true,
		Response: state,
	})
}

func err500(w http.ResponseWriter, err error) {
	log.Println("err 500: %w", err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(Response[Empty]{
		OK:    false,
		Error: err.Error(),
	})
}

type StateMessage StreamState

func (h StateMessage) Type() string { return "state" }

func (a *API) StateWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsrpc.NewConn(w, r, wsrpc.CommandPalette{})
	if err != nil {
		err500(w, err)
		return
	}
	defer conn.Close()

	pump := conn.Pump(r.Context())

	state, err := a.fetchState()
	if err != nil {
		err500(w, err)
		return
	}
	conn.SendMessage(StateMessage(state))

	sub, err := a.stateUpdates.Subscribe(r.Context())
	if err != nil {
		err500(w, err)
		return
	}
	defer sub.Unsubscribe(context.Background())

	for {
		select {
		case update := <-sub.Ch():
			conn.SendMessage(StateMessage(update))
		case _, more := <-pump.Ch():
			if !more {
				// Client disconnected.
				return
			}
		}
	}
}
