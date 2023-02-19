package chat

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/irth/broadcast"
	"github.com/irth/wsrpc"
)

var DEBUG = os.Getenv("DEBUG") == "1"

type Chat struct {
	ch *broadcast.Channel[wsrpc.Message]

	nicks     map[string]struct{}
	nicksLock sync.Mutex
}

func NewChat() (*Chat, error) {
	return &Chat{
		ch: broadcast.NewChannel[wsrpc.Message](),

		nicks: make(map[string]struct{}),
	}, nil
}

func (c *Chat) Handler() http.Handler {
	r := chi.NewRouter()
	r.HandleFunc("/ws", c.WebsocketHandler)
	return r
}

func (c *Chat) Start(ctx context.Context) {
	go c.ch.Run(ctx)
}

func (c *Chat) Listen(addr string) error {
	return http.ListenAndServe(addr, c.Handler())
}
