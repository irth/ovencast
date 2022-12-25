package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a *API) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from ovencast api :)\n")
	})
	return r
}
