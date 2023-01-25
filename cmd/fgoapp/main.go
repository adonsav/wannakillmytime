package main

import (
	"fmt"
	"github.com/adonsav/fgoapp/pkg/config"
	"github.com/adonsav/fgoapp/pkg/handlers"
	"github.com/adonsav/fgoapp/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var (
	appConfig config.AppConfig
	session   *scs.SessionManager
)

// Entry point
func main() {
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot crate template cache")
	}
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	// alternatively we can use the below in place of the two method calls
	// above and  delete NewRepo and NewHandlers methods
	// handlers.Repo = &handlers.Repository{App: &appConfig}

	render.NewTemplates(&appConfig)

	fmt.Printf("Starting application on port %s", portNumber)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	log.Fatal(server.ListenAndServe())
}
