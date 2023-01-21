package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gobwas/ws"
)

type Ping = WSCommand[PingCommand, PingResponse]

type PingCommand struct {
	Ping string `json:"ping"`
}

type PingResponse struct {
	Pong string `json:"pong"`
}

func (a *API) Websocket(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(Response[Empty]{
			OK:    false,
			Error: err.Error(),
		})
	}

	wsconn := NewWSConn(conn)

	for {
		cmd, err := wsconn.DecodeCommand()
		if err != nil {
			log.Printf("ws: decode command: %s", err)
			return
		}

		switch cmd := cmd.(type) {
		case Ping:
			cmd.Reply(PingResponse{
				Pong: cmd.Request.Ping,
			})
		default:
			cmd.Error("unimplemented command")
		}
	}

}
