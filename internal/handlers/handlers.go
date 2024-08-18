package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andreadebortoli2/GO-bnb/internal/config"
	"github.com/andreadebortoli2/GO-bnb/internal/forms"
	"github.com/andreadebortoli2/GO-bnb/internal/models"
	"github.com/andreadebortoli2/GO-bnb/internal/render"
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

	render.RenderTemplates(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	// pass some data to about page
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, passing data successfully"

	// take the IP address from the session and send to template
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplates(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals renders the genereal's quarters room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Reservation renders the major's suite room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Reservation renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handle request for availability and send json
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	response := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	responseMarshaled, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseMarshaled)
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// use validation functions
	// form.Has("first-name", r)
	form.Required("first-name", "last-name", "email")
	form.MinLength("first-name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplates(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// put reservation in session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// redirect after post
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get reservation from session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// if there is no reservation in session cause of error or user bad behavior send back to home page with error message
		log.Println("cannot get item from the session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// remove reservation from session
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplates(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
