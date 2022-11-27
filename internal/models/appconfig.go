package models

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"net/http"
)

type AppConfig struct {
	Addr           string
	DSN            string
	JWTSecret      string
	APIKey         string
	DB             *sql.DB
	TemplateCache  map[string]*template.Template
	IsProduction   bool
	Username       string
	Password       string
	SessionManager *scs.SessionManager
	IsLoggedIn     bool
}

func (app *AppConfig) IsAuthenticated(r *http.Request) bool {
	return app.SessionManager.Exists(r.Context(), "authenticatedUserID")
}
