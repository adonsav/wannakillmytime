package main

import (
	"encoding/gob"
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/handlers"
	"github.com/adonsav/fgoapp/internal/helpers"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/adonsav/fgoapp/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var (
	appConfig config.AppConfig
	session   *scs.SessionManager
	infoLog   *log.Logger
	errorLog  *log.Logger
)

// Entry point
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port%s\n", portNumber)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	log.Fatal(server.ListenAndServe())
}

func run() error {
	gob.Register(models.Registration{})
	appConfig.InProduction = false

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot crate template cache")
		return err
	}
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	// alternatively we can use the below in place of the two method calls
	// above and  delete NewRepo and NewHandlers methods
	// handlers.Repo = &handlers.Repository{handlersAppConfig: &appConfig}
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)

	return nil
}
