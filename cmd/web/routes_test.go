package main

import (
	"testing"

	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing, test passed
	default:
		t.Errorf("type is not an http.handler, but is %T", v)
	}
}
