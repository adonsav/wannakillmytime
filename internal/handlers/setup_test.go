package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/adonsav/fgoapp/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var handlersTestAppConfig config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../templates/gohtmltemplates"

// functions allows us to specify certain functions available to Go templates
var functions = template.FuncMap{
	"humanDate": render.HumanDate,
}

func TestMain(m *testing.M) {
	handlersTestAppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	handlersTestAppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	gob.Register(models.User{})
	gob.Register(models.Registration{})
	gob.Register(models.BoredBuddie{})
	handlersTestAppConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = handlersTestAppConfig.InProduction

	handlersTestAppConfig.Session = session

	emailChan := make(chan models.EmailData)
	handlersTestAppConfig.EmailChan = emailChan
	defer close(emailChan)
	listenForEmail()

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	handlersTestAppConfig.TemplateCache = templateCache
	handlersTestAppConfig.UseCache = true

	repo := NewTestRepo(&handlersTestAppConfig)
	NewHandlers(repo)
	// alternatively we can use the below in place of the two method calls
	// above and  delete NewRepo and NewHandlers methods
	// handlers.Repo = &handlers.Repository{handlersAppConfig: &handlersTestAppConfig}
	render.NewRenderer(&handlersTestAppConfig)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/register", Repo.Registration)
	mux.Post("/register", Repo.PostRegistration)
	mux.Get("/registration-summary", Repo.RegistrationSummary)

	mux.Get("/question-mark", Repo.QuestionMark)
	mux.Post("/question-mark-json", Repo.QuestionMarkJSON)

	mux.Get("/user/login", Repo.Login)
	mux.Post("/user/login", Repo.PostLogin)
	mux.Get("/user/logout", Repo.Logout)

	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/registrations-active", Repo.AdminActiveRegistrations)
	mux.Get("/admin/registrations-all", Repo.AdminAllRegistrations)
	mux.Get("/admin/registrations/{src}/{id}", Repo.AdminShowRegistration)
	mux.Get("/admin/deactivate-registration/{src}/{id}", Repo.AdminDeactivateRegistration)
	mux.Post("/admin/registrations/{src}/{id}", Repo.AdminPostShowRegistration)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// noSurf adds CSRF protection to all POST requests
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   handlersTestAppConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session in every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files ending with *.page.gohtml from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.gohtml
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}

// listenForEmail creates a goroutine that just listens for emails in the respective channels without
// taking further actions
func listenForEmail() {
	go func() {
		for {
			_ = <-handlersTestAppConfig.EmailChan
		}
	}()
}
