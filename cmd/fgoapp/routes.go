package main

import (
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(noSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/register", handlers.Repo.Registration)
	mux.Post("/register", handlers.Repo.PostRegistration)
	mux.Get("/registration-summary", handlers.Repo.RegistrationSummary)

	mux.Get("/question-mark", handlers.Repo.QuestionMark)
	mux.Post("/question-mark-json", handlers.Repo.QuestionMarkJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
