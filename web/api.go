package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/irth/broadcast"
	"github.com/irth/ovencast/web/ome"
)

// API implements the http/websocket API that is used by OME for access control
// and by the web frontend for everything else, including stream state updates
// and chat.
type API struct {
	*Config

	OME *ome.API

	// stream state
	state                  StreamState
	stateUpdates           *broadcast.Channel[StreamState]
	admissionWebhookSignal chan bool

	// chat
	nicks    map[string]struct{}
	nickLock sync.Mutex
}

// NewAPI loads config from the specified path and instantiates the API struct.
// You still need to start its main loop by calling Run.
func NewAPI(configPath string) (*API, error) {
	conf, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	api, err := ome.NewAPI(conf.OME.APIAddr, conf.OME.APIToken)
	if err != nil {
		return nil, fmt.Errorf("NewOMEAPI: %w", err)
	}

	return &API{
		Config: conf,
		OME:    api,

		stateUpdates:           broadcast.NewChannel[StreamState](),
		admissionWebhookSignal: make(chan bool, 8),

		nicks: make(map[string]struct{}),
	}, nil
}

// Run starts the API handler main loop. This pretty much never exits, so start
// it in a goroutine if you need to.
func (a *API) Run() {
	// start the broadcast channel's event loop
	go a.stateUpdates.Run(context.TODO())

	// periodically check for state updates
	go a.runStateUpdater()

	stateSub, _ := a.stateUpdates.Subscribe(context.TODO())
	for state := range stateSub {
		log.Printf("state updated! %+v", state)
	}
}

// Router creates the http.Handler with API HTTP routes configured
func (a *API) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "hello from ovencast api :)\n")
	})
	r.Get("/state", a.State)

	r.Get("/websocket", a.Websocket)

	r.Post("/admission", a.AdmissionWebhook)
	return r
}

type Empty any

// Response defines the structure of the HTTP API responses. All endpoints
// should use it.
type Response[T any] struct {
	OK       bool   `json:"ok"`
	Error    string `json:"error,omitempty"`
	Response T      `json:"response,omitempty"`
}

// State (GET <api>/state) fetches the up-to-date stream state from OME and
// returns it to the client.
func (a *API) State(w http.ResponseWriter, r *http.Request) {
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
