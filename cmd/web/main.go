package main

import (
	"encoding/gob"
	"flag"
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

	// declare what i'm going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// CLI flags (type pointers)
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPassword := flag.String("dbpassword", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("missing required flags")
		os.Exit(1)
	}

	// open a channel to serve emails through, defer in main not in run so the channel will not be closed as soon as the functionrun is executed
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

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
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPassword, *dbSSL)
	// dev mode connectionString: "host=localhost port=5432 dbname=go_bnb user=postgres password=password"
	db, err := driver.ConnectSQL(connectionString)
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
