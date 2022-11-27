package main

import (
	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/elkcityhazard/go-andrew-mccall/internal/render"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
	"github.com/justinas/alice"
)

func routes() http.Handler {

	router := httprouter.New()

	dynamic := alice.New(app.SessionManager.LoadAndSave, noSurf)

	protected := dynamic.Append(RequireAuthentication)

	router.NotFound = dynamic.ThenFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.SessionManager.Put(r.Context(), "flash", "I am sorry, but that page can not be found")

		render.RenderTemplate(w, r, "404.tmpl.html", &models.TemplateData{})
	}))

	fs := http.FileServer(http.Dir("./static"))

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fs))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(handlers.Repo.Home))
	router.Handler(http.MethodGet, "/posts/", dynamic.ThenFunc(handlers.Repo.GetListOfPosts))
	router.Handler(http.MethodGet, "/posts/:id", dynamic.ThenFunc(handlers.Repo.GetSinglePost))

	router.Handler(http.MethodGet, "/admin/signup", dynamic.ThenFunc(handlers.Repo.Signup))
	router.Handler(http.MethodGet, "/admin/login", dynamic.ThenFunc(handlers.Repo.Login))
	router.Handler(http.MethodPost, "/admin/login", dynamic.ThenFunc(handlers.Repo.Login))
	router.Handler(http.MethodGet, "/admin/logout", dynamic.ThenFunc(handlers.Repo.Logout))
	router.Handler(http.MethodGet, "/admin/add-post", protected.ThenFunc(handlers.Repo.AddPost))
	router.Handler(http.MethodPost, "/admin/add-post", protected.ThenFunc(handlers.Repo.AddPost))
	router.Handler(http.MethodPost, "/admin/get-jwt", dynamic.ThenFunc(handlers.Repo.GetJWT))

	standard := alice.New(secureHeaders)

	return standard.Then(router)

}
