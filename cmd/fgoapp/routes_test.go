package main

import (
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch h := mux.(type) {
	case *chi.Mux:
		// we are good to go, do nothing
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but is %T", h))
	}
}
