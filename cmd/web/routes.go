package main

import (
	"net/http"

	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/andreadebortoli2/GO-bnb/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes handle all the routes
func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	// recover from panics, log the panic and return an HTTP 500 status(if possible)
	mux.Use(middleware.Recoverer)
	// add crsf token in cookies for POST protection
	mux.Use(NoSurf)
	// load and use the session
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)

	// define file server where to take static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
