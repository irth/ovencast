package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/irth/wsrpc"
)

var palette = wsrpc.CommandPalette{
	"nick":    NickCommand{},
	"message": MessageCommand{},
}

type ClientState struct {
	ctx context.Context

	ready bool
	nick  string
}

func err500(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s: %s", msg, err.Error())
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(map[string]any{
		"ok":    false,
		"error": fmt.Sprintf("%s: %s", msg, err.Error()),
	})
}

func (c *Chat) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade to websocket
	conn, err := wsrpc.NewConn(w, r, palette)
	if err != nil {
		err500(w, "websocket upgrade failed", err)
		return
	}
	defer conn.Close()

	// Subscribe to broadcasts
	broadcasts, err := c.ch.Subscribe(r.Context())
	if err != nil {
		err500(w, "broadcast channel sub failed", err)
		return
	}
	// TODO: implement unsubs in broadcast

	// Set up client state
	state := ClientState{
		ctx:   r.Context(),
		ready: false,
	}

	// Let go of the nick if the client disconnects
	defer func() {
		if state.ready {
			c.nicksLock.Lock()
			defer c.nicksLock.Unlock()
			delete(c.nicks, state.nick)
		}
	}()

	// Read the websocket messages into a channel so that we can select{} it
	cmds := conn.Pump(r.Context())

	// Aaaand the main event loop
	for {
		select {
		// Return in case the client disconnects (request context expires)
		case <-r.Context().Done():
			log.Println("request context expired, exiting")
			return

		// Push broadcasts to clients
		case broadcast, more := <-broadcasts:
			if !more {
				return
			}
			err := conn.SendMessage(broadcast)
			if err != nil {
				log.Printf("failed to send message to websocket: %s", err.Error())
				return
			}

		// Listen to client commands
		case cmd, more := <-cmds.Ch():
			if !more {
				if err := cmds.Err(); err != nil {
					log.Printf("chat: websocket decode error: %s", err.Error())
					return
				}
			}

			if DEBUG {
				log.Printf("got command: %+v", cmd)
				log.Printf("state before: %+v", state)
			}

			switch cmd := cmd.(type) {
			case NickCommand:
				c.handleNick(&state, cmd)
			case MessageCommand:
				c.handleMessage(&state, cmd)

			default:
				cmd.Err("not yet implemented")
			}
			if DEBUG {
				log.Printf("state after: %+v", state)
			}
		}

	}
}
