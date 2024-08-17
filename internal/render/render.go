package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/andreadebortoli2/GO-bnb/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates set config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// ? prepared for later usage
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplates renders templates using html/template
func RenderTemplates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var templateCache map[string]*template.Template

	if app.UseCache {
		// production mode (instead of create the template cache each time, i want to use the app config template cache)
		templateCache = app.TemplateCache
	} else {
		// dev mode (re-create the template cache each time to check changes i make developing)
		templateCache, _ = CreateTemplateCache()
	}

	requestedTemplate, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	// ? preapared for later usage
	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)
	err := requestedTemplate.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = templateSet
	}

	return myCache, nil
}
