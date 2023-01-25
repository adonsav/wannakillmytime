package handlers

import (
	"github.com/adonsav/fgoapp/pkg/config"
	"github.com/adonsav/fgoapp/pkg/models"
	"github.com/adonsav/fgoapp/pkg/render"
	"net/http"
)

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(config *config.AppConfig) *Repository {
	return &Repository{
		App: config,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the home page handler
func (hr *Repository) Home(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Test"
	render.Template(w, "home.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// About is the about page handler
func (hr *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "about.page.gohtml", &models.TemplateData{})
}
