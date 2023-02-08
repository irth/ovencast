package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/irth/ovencast/web/ws"
)

type PingCommand struct {
	Ping string `json:"ping"`
}

type PingResponse struct {
	Pong string `json:"pong"`
}

type Ping = ws.Command[PingCommand, PingResponse]

var websocketCommands = ws.CommandPalette{
	"ping": Ping{},
}

type HelloMessage struct {
	Version string `json:"version"`
}

func (h HelloMessage) Type() string { return "hello" }

type StateMessage StreamState

func (h StateMessage) Type() string { return "state" }

func err500(w http.ResponseWriter, err error) {
	log.Println("err 500: %w", err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(Response[Empty]{
		OK:    false,
		Error: err.Error(),
	})
}

func (a *API) Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.NewConn(w, r, websocketCommands)
	if err != nil {
		err500(w, err)
		return
	}
	defer conn.Close()

	ctx := r.Context()

	stateUpdates, err := a.stateUpdates.Subscribe(ctx)
	// TODO: unsubscribe!!!
	if err != nil {
		err500(w, err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case state := <-stateUpdates:
				conn.SendMessage(StateMessage(state))
			}
		}
	}()

	err = conn.SendMessage(HelloMessage{
		// TODO: do some build magic to have git version numbers
		Version: "0.0.1",
	})

	if err != nil {
		err500(w, err)
		return
	}

	state, err := a.fetchState()
	if err != nil {
		err500(w, err)
		return
	}

	err = conn.SendMessage(StateMessage(state))
	if err != nil {
		err500(w, err)
		return
	}

	for {
		cmd, err := conn.Decode()
		if err != nil {
			log.Printf("ws: decode command: %s", err)
			return
		}

		switch cmd := cmd.(type) {
		case Ping:
			cmd.OK(PingResponse{
				Pong: cmd.Request.Ping,
			})
		default:
			cmd.Err("unimplemented command")
		}
	}
}
