package main

import (
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.Repo.Home)

	mux.HandleFunc("/admin/signup", handlers.Repo.Signup)
	mux.HandleFunc("/admin/login", handlers.Repo.Login)
	mux.Handle("/admin/add-post", utils.ValidateJWT(handlers.Repo.AddPost))
	mux.HandleFunc("/admin/get-jwt", handlers.Repo.GetJWT)

	return mux
}
