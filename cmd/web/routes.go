package main

import (
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.Repo.Home)

	return mux
}
