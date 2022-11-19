package handlers

import (
	"fmt"
	"net/http"

	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/elkcityhazard/go-andrew-mccall/internal/render"
)

type Repository struct {
	AppConfig *models.AppConfig
}

var app *models.AppConfig

var Repo *Repository

func NewHandlers(a *models.AppConfig) {
	app = a
}

func NewRepo(app *models.AppConfig) *Repository {
	return &Repository{
		AppConfig: app,
	}
}

func SetRepo(m *Repository) {
	Repo = m
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	pathKey := r.URL.Path[len("/"):]

	fmt.Println(pathKey)

	fmt.Println("is this working?")

	render.RenderTemplate(w, r, "home.tmpl.html", &models.TemplateData{})
}
