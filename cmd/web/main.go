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
	"github.com/andreadebortoli2/GO-bnb/internal/driver"
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

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	fmt.Println("Starting mail listener...")
	listenForMail()
	// test the mail channel
	/* msg := models.MailData{
		To:      "john@doe.ca",
		From:    "me@here.com",
		Subject: "Some subject",
		Content: "",
	}
	app.MailChan <- msg */

	// automatically send email at the start with standard library
	/* from := "me@here.com"
	auth:= smtp.PlainAuth("",from,"","localhost")
	err = smtp.SendMail("localhost:1025",auth,from,[]string["you@there.com"],[]byte("Hello, world"))
	if err != nil {
		log.Println(err)
	} */

	// set the server
	_, _ = fmt.Printf("Starting application on port %s \n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	// what i'm going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// open a channel to serve emails through, defer in main not in run so the channel will not be closed as soon as the functionrun is executed
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

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

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go_bnb user=postgres password=password")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database")

	// generate template cache and save in AppConfig
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = templateCache

	// set the config + database connection as repo for the handlers and set the config + database available to handlers pkg
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// set the config available to render pkg
	render.NewRenderer(&app)

	// set the config available to helpers pkg
	helpers.NewHelpers(&app)

	return db, nil
}
