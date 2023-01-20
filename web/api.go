package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/irth/ovencast/web/ome"
)

type API struct {
	*Config

	OME *ome.API
}

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
	}, nil
}

func (a *API) Run() {
	// TODO: do stuff
	for {
		time.Sleep(time.Second * 1)
	}
}

func (a *API) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from ovencast api :)\n")
	})
	r.Get("/online", a.Online)
	r.Post("/admission", a.AdmissionWebhook)
	return r
}

type Empty interface{}

type Response[T any] struct {
	OK       bool   `json:"ok"`
	Error    string `json:"error,omitempty"`
	Response T      `json:"response,omitempty"`
}

type StreamOnlineResponse struct {
	Online bool `json:"online"`
}

func (a *API) Online(w http.ResponseWriter, r *http.Request) {
	online, err := a.OME.StreamExists("default", "live", "stream")

	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response[Empty]{
			OK:    false,
			Error: err.Error(),
		})
	}

	json.NewEncoder(w).Encode(Response[StreamOnlineResponse]{
		OK:       true,
		Response: StreamOnlineResponse{Online: online},
	})
}
