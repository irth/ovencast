package chat

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/irth/broadcast"
)

type Chat struct {
	ch *broadcast.Channel[any] // TODO: create a type for messages

	nicks     map[string]struct{}
	nicksLock sync.Mutex
}

func NewChat() (*Chat, error) {
	return &Chat{
		ch: broadcast.NewChannel[any](),

		nicks: make(map[string]struct{}),
	}, nil
}

func (c *Chat) Listen(addr string) error {
	r := chi.NewRouter()
	r.HandleFunc("/ws", c.WebsocketHandler)

	return http.ListenAndServe(addr, r)
}
