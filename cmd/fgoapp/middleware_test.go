package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var tHandler testingHandler
	handler := noSurf(&tHandler)

	switch h := handler.(type) {
	case http.Handler:
		// we are good to go, do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", h))
	}
}

func TestSessionLoad(t *testing.T) {
	var tHandler testingHandler
	handler := SessionLoad(&tHandler)

	switch h := handler.(type) {
	case http.Handler:
		// we are good to go, do nothing
	default:
		t.Error(fmt.Sprintf("type us not an http.Handler, but is %T", h))
	}
}
