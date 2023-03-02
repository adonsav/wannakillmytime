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
	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		//mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/registrations-active", handlers.Repo.AdminActiveRegistrations)
		mux.Get("/registrations-all", handlers.Repo.AdminAllRegistrations)
		mux.Get("/registrations/{src}/{id}", handlers.Repo.AdminShowRegistration)
		mux.Get("/deactivate-registration/{src}/{id}", handlers.Repo.AdminDeactivateRegistration)
		mux.Post("/registrations/{src}/{id}", handlers.Repo.AdminPostShowRegistration)
	})
	return mux
}
