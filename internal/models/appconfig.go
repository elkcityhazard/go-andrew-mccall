package models

import (
	"database/sql"
	"html/template"
)

type AppConfig struct {
	Addr          string
	DSN           string
	JWTSecret     string
	APIKey        string
	DB            *sql.DB
	TemplateCache map[string]*template.Template
	IsProduction  bool
	Username      string
	Password      string
}
