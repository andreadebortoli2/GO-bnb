package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/andreadebortoli2/GO-bnb/internal/handlers"
	"github.com/andreadebortoli2/GO-bnb/internal/helpers"
	"github.com/andreadebortoli2/GO-bnb/internal/models"
	"github.com/andreadebortoli2/GO-bnb/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}

	// set the server
	_, _ = fmt.Printf("Starting application on port %s \n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	// what i'm going to store in the session
	gob.Register(models.Reservation{})

	// change to true when in production
	app.InProduction = false
	app.UseCache = false

	// set info and error logs
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// define the session and save in AppConfig
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// generate template cache and save in AppConfig
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}
	app.TemplateCache = templateCache

	// set the config as repo for the handlers and set the config available to handlers pkg
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// set the config available to render pkg
	render.NewTemplates(&app)

	// set the config available to helpers pkg
	helpers.NewHelpers(&app)

	return nil
}
