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

	flag.StringVar(&app.Username, "username", "nil", "username")
	flag.StringVar(&app.Password, "password", "nil", "password")
	flag.StringVar(&app.JWTSecret, "jwtsecret", "nil", "JWT Secret")
	flag.StringVar(&app.APIKey, "apikey", "", "apikey to enable auth")

	flag.Parse()

	fmt.Println("password", app.Password)

	fmt.Println(app)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(192.168.0.4:3306)/andrew_mccall", app.Username, app.Password))

	if err != nil {
		log.Fatalln(err)
	}

	// err = db.Ping()

	// if err != nil {
	// 	fmt.Println(err)
	// 	log.Fatalln(err)
	// }

	defer db.Close()

	app.Addr = ":8080"
	app.DB = db
	app.IsProduction = false

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = tc

	//	New handlers

	handlers.NewHandlers(&app)

	// Set Up Repos
	repo := handlers.NewRepo(&app)
	handlers.SetRepo(repo)

	// Set Up A New Renderer And Give It Access To AppConfig=

	render.NewRenderer(&app)

	// New Utils

	utils.NewUtils(&app)

	if err != nil {
		log.Fatalln(err)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("listening on: ", srv.Addr)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
