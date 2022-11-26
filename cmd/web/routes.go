package main

import (
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.Repo.Home)
	mux.HandleFunc("/posts/", handlers.Repo.GetListOfPosts)

	mux.HandleFunc("/admin/signup", handlers.Repo.Signup)
	mux.HandleFunc("/admin/login", handlers.Repo.Login)
	mux.HandleFunc("/admin/add-post", handlers.Repo.AddPost)
	mux.Handle("/admin/get-jwt", IsLoggedIn(http.HandlerFunc(handlers.Repo.GetJWT)))

	return CheckForAPIKey(SetAPIKey(mux))
}
