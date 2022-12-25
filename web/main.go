package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var static embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatalf("loading static files: %s", err)
	}
	staticHandler := http.FileServer(http.FS(staticFS))

	api := NewAPI()
	go api.Run()

	r.Mount("/api", api.Router())
	r.Mount("/", staticHandler)

	http.ListenAndServe(":8080", r)
}
