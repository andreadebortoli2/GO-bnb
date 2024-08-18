package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf("type is not an http.handler, but is %T", v)
	}
}
func TestSessioLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf("type is not an http.handler, but is %T", v)
	}
}
