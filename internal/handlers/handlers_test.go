package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "get", []postData{}, http.StatusOK},
	{"about", "/about", "get", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "get", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "get", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "get", []postData{}, http.StatusOK},
	{"contact", "/contact", "get", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", "post", []postData{
		{key: "start", value: "2024-01-01"},
		{key: "end", value: "2024-01-02"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", "post", []postData{
		{key: "start", value: "2024-01-01"},
		{key: "end", value: "2024-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "post", []postData{
		{key: "first-name", value: "John"},
		{key: "last-name", value: "Doe"},
		{key: "email", value: "john@doe.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "get" {
			resp, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
