package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreadebortoli2/GO-bnb/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "get", http.StatusOK},
	{"about", "/about", "get", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "get", http.StatusOK},
	{"majors-suite", "/majors-suite", "get", http.StatusOK},
	{"search-availability", "/search-availability", "get", http.StatusOK},
	{"contact", "/contact", "get", http.StatusOK},

	// {"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},
	// {"post-search-availability", "/search-availability", "post", []postData{
	// 	{key: "start", value: "2024-01-01"},
	// 	{key: "end", value: "2024-01-02"},
	// }, http.StatusOK},
	// {"post-search-availability-json", "/search-availability-json", "post", []postData{
	// 	{key: "start", value: "2024-01-01"},
	// 	{key: "end", value: "2024-01-02"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		resp, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	// prepare data to create a session for testing
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	// get the reservation func as a handler to serve directly without the routes
	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {

	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test with missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler failed when trying to fail inserting reservation into database, returned code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for failure to insert room restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first-name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last-name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@example.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler failed when trying to fail inserting room restriction into database, returned code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
