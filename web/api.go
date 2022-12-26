package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type API struct {
	*Config
}

func NewAPI(configPath string) (*API, error) {
	conf, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &API{
		Config: conf,
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
	r.Post("/admission", a.AdmissionWebhook)
	return r
}
