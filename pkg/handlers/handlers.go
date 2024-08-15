package handlers

import (
	"net/http"

	"github.com/andreadebortoli2/GO-bnb/pkg/config"
	"github.com/andreadebortoli2/GO-bnb/pkg/models"
	"github.com/andreadebortoli2/GO-bnb/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// add to the session the IP address of the browser hitting the home
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplates(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// pass some data to about page
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, passing data successfully"

	// take the IP address from the session and send to template
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplates(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// Generals renders the genereal's quarters room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "generals.page.tmpl", &models.TemplateData{})
}

// Reservation renders the major's suite room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "majors.page.tmpl", &models.TemplateData{})
}

// Reservation renders the major's suite room page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "contact.page.tmpl", &models.TemplateData{})
}
