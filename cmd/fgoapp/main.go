package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(appConfig.EmailChan)

	fmt.Println("Starting email listener...")
	listenForEmail()

	fmt.Printf("Starting application on port%s\n", portNumber)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	log.Fatal(server.ListenAndServe())
}

func run() (*driver.DB, error) {
	gob.Register(models.User{})
	gob.Register(models.Registration{})
	gob.Register(models.BoredBuddie{})

	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	appConfig.InProduction = *inProduction

	emailChan := make(chan models.EmailData)
	appConfig.EmailChan = emailChan

	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		*dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Adios...")
	}
	log.Println("Connected to database")

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = *useCache

	repo := handlers.NewRepo(&appConfig, db)
	handlers.NewHandlers(repo)
	// alternatively we can use the below in place of the two method calls
	// above and  delete NewRepo and NewHandlers methods
	// handlers.Repo = &handlers.Repository{handlersAppConfig: &appConfig}
	render.NewRenderer(&appConfig)
	helpers.NewHelpers(&appConfig)

	return db, nil
}
