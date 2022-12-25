package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type API struct {
	ControlService string
}

func NewAPI() *API {
	return &API{ControlService: "http://control:9595"} // TODO: allow this to be configurable
}

func (a *API) Run() {
	log.Println("Waiting for the control service...")

	tries := 0
	maxTries := 10
	for {
		r, err := http.Get(a.ControlService + "/ping")
		tries = tries + 1
		if err == nil && r.StatusCode == http.StatusOK {
			log.Println("Control service up.")
			break
		}

		triesStr := fmt.Sprintf("(%d/%d tries)", tries, maxTries)

		if err != nil {
			log.Printf("Failed to connect to the control service %s: %s", triesStr, err)
		} else {
			log.Printf("Failed to connect to the control service %s: http status %s", triesStr, r.Status)
		}

		if tries == maxTries {
			log.Fatal("No tries left - giving up.")
		}

		log.Printf("Retrying in 5 seconds ...")
		time.Sleep(time.Second * 5)
	}

	// TODO: set stream key
	for {
		time.Sleep(time.Second * 1)
	}
}

func (a *API) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from ovencast api :)\n")
	})
	return r
}
