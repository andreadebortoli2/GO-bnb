package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/andreadebortoli2/GO-bnb/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var app *config.AppConfig

var pathToTemplates = "./templates"

// NewRenderer set config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate return time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate format the date in strings
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate returns a slice of ints, starting at 1, going to count
func Iterate(count int) []int {
	var items []int
	for i := 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

// Add two ints
func Add(a, b int) int {
	return a + b
}

// AddDefaultData add data default/placeholder to TemplateData
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Templates renders templates using html/template
func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

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
		return errors.New("can't get template from template cache")
	}

	// default/placeholder data from models.TemplateData
	td = AddDefaultData(td, r)

	buf := new(bytes.Buffer)
	err := requestedTemplate.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = templateSet
	}

	return myCache, nil
}
