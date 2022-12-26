package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var static embed.FS

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

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("loading static files: %s", err)
	}
	staticHandler := http.FileServer(http.FS(staticFS))

	api, err := NewAPI(configPath)
	if err != nil {
		log.Fatalf("creating api service: %s", err)
	}
	go api.Run()

	r.Mount("/api", api.Router())
	r.Mount("/", staticHandler)

	http.ListenAndServe(listenAddr, r)
}
