package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	listenAddr, ok := os.LookupEnv("OVENCAST_WEB_ADDR")
	if !ok {
		listenAddr = ":8080"
	}

	configPath, ok := os.LookupEnv("OVENCAST_WEB_CONF")
	if !ok {
		configPath = "./config.yaml"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	staticFS := http.Dir("./static")
	staticHandler := http.FileServer(staticFS)

	api, err := NewAPI(configPath)
	if err != nil {
		log.Fatalf("creating api service: %s", err)
	}
	go api.Run()

	r.Mount("/api", api.Router())
	r.Mount("/", staticHandler)

	http.ListenAndServe(listenAddr, r)
}
