package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/irth/wsrpc"
)

var palette = wsrpc.CommandPalette{
	"nick": NickCommand{},
}

type NickCommand = wsrpc.Command[string, any]

type ClientState struct {
	ready bool
	nick  string
}

func (c *Chat) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsrpc.NewConn(w, r, palette)
	if err != nil {
		log.Printf("websocket upgrade failed: %s", err.Error())
		json.NewEncoder(w).Encode(map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	state := ClientState{
		ready: false,
	}
	defer func() {
		if state.ready {
			// means we have a nick, so we gotta let go
			c.nicksLock.Lock()
			defer c.nicksLock.Unlock()
			delete(c.nicks, state.nick)
		}
	}()

	for {
		cmd, err := conn.Decode()
		if err != nil {
			log.Printf("websocket decode error: %s", err.Error())
			return
		}

		switch cmd := cmd.(type) {
		case NickCommand:
			c.handleNick(&state, cmd)

		default:
			cmd.Err("not yet implemented")
		}
	}
}
