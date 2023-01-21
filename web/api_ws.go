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

var websocketCommands = ws.CommandPallete{
	"ping": Ping{},
}

type HelloMessage struct {
	Version string
	Online  bool
}

func (h HelloMessage) Type() string { return "hello" }

func (a *API) Websocket(w http.ResponseWriter, r *http.Request) {
	defer log.Println("dupa")
	conn, err := ws.NewConn(w, r, websocketCommands)
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response[Empty]{
			OK:    false,
			Error: err.Error(),
		})
		return
	}
	defer conn.Close()

	online, _ := a.isOnline()

	err = conn.SendMessage(HelloMessage{
		Version: "0.0.1",
		Online:  online,
	})
	if err != nil {
		log.Printf("ws sendmessage: %s", err)
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
