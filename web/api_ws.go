package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/irth/wsrpc"
)

type PingCommand struct {
	Ping string `json:"ping"`
}

type PingResponse struct {
	Pong string `json:"pong"`
}

type Ping = wsrpc.Command[PingCommand, PingResponse]

type NickCommand struct {
	Nickname string `json:"nickname"`
}

type NickResponse any // always nil

type Nick = wsrpc.Command[NickCommand, NickResponse]

var websocketCommands = wsrpc.CommandPalette{
	"ping": Ping{},
	"nick": Nick{},
}

type HelloMessage struct {
	Version string `json:"version"`
}

func (h HelloMessage) Type() string { return "hello" }

type StateMessage StreamState

func (h StateMessage) Type() string { return "state" }

type NickChangeMessage struct {
	Previous string `json:"previous"`
	New      string `json:"new"`
}

func (n NickChangeMessage) Type() string { return "nick" }

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
	conn, err := wsrpc.NewConn(w, r, websocketCommands)
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

	nick := ""
	nickChangeThrottle := time.Now()

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

		case Nick:
			if nick == cmd.Request.Nickname {
				cmd.OK(nil)
				continue
			}

			if valid := validateNick(cmd.Request.Nickname); !valid {
				cmd.Err("nickname invalid") // TODO: error codes for the frontend?
				continue
			}

			if nickChangeThrottle.After(time.Now()) {
				cmd.Err("too many nick change requests, please wait")
				continue
			}

			a.nickLock.Lock()
			_, inUse := a.nicks[cmd.Request.Nickname]
			if !inUse {
				delete(a.nicks, nick)
				a.nicks[cmd.Request.Nickname] = struct{}{}
				// TODO: clear nick on disconnect
			}
			a.nickLock.Unlock()

			if inUse {
				cmd.Err("nickname in use")
				continue
			}

			log.Printf("nickname change: %s -> %s", nick, cmd.Request.Nickname)

			nick = cmd.Request.Nickname
			nickChangeThrottle = time.Now().Add(30 * time.Second)

			cmd.OK(nil)
			

		default:
			cmd.Err("unimplemented command")
		}
	}
}

// TODO: support something more fun than ascii only
var NickRegex = regexp.MustCompile(`^[a-zA-Z0-9\._-]{3,32}$`)

func validateNick(n string) bool {
	return NickRegex.MatchString(n)
}
