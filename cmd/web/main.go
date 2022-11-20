package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elkcityhazard/go-andrew-mccall/internal/handlers"
	"github.com/elkcityhazard/go-andrew-mccall/internal/models"
	"github.com/elkcityhazard/go-andrew-mccall/internal/render"
	"github.com/elkcityhazard/go-andrew-mccall/internal/utils"
	_ "github.com/go-sql-driver/mysql"
)

var app models.AppConfig

func main() {

	flag.StringVar(&app.DSN, "dsn", "", "the data source name string")
	flag.StringVar(&app.JWTSecret, "jwtsecret", "", "JWT Secret")
	flag.StringVar(&app.APIKey, "apikey", "", "apikey to enable auth")

	flag.Parse()

	db, err := sql.Open("mysql", app.DSN)

	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	app.Addr = ":8080"
	app.DB = db
	app.IsProduction = false

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = tc

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	handlers.NewHandlers(&app)

	// Set Up Repos
	repo := handlers.NewRepo(&app)
	handlers.SetRepo(repo)

	// Set Up A New Renderer And Give It Access To AppConfig=

	render.NewRenderer(&app)

	// New Utils

	utils.NewUtils(&app)

	fmt.Println("listening on: ", srv.Addr)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
