package main

import (
	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
	"net/http"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Repo.Home)

	return mux
}
