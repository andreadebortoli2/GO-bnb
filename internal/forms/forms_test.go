package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("get invalid with no errors")
	}

	form.Errors.Add("test error", "test error")
	isValidWError := form.Valid()
	if isValidWError {
		t.Error("get valid with errors")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required filds when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)

	form.Has("a")
	if form.Valid() {
		t.Error("show valid when required field missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Has("a")
	if !form.Valid() {
		t.Error("show not valid when required field present")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("a", 3)
	if form.Valid() {
		t.Error("show valid when the field value string is not longh enough")
	}

	postedData := url.Values{}
	postedData.Add("a", "test")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.MinLength("a", 10)
	if form.Valid() {
		t.Error("show valid when the field value string is not long enough")
	}

	// test also forms.errors.get()
	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "test")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("show not valid when the field value string is long enough")
	}

	// test also forms.errors.get()
	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error, but did get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.PostForm)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("show valid when email is not a valid email")
	}

	postedData := url.Values{}
	postedData.Add("a", "aa@aa.a")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.IsEmail("a")
	if !form.Valid() {
		t.Error("show not valid when email is correct")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	r, _ = http.NewRequest("post", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.IsEmail("a")
	if form.Valid() {
		t.Error("show valid when email is not correct")
	}
}
