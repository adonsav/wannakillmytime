package handlers

import (
	"encoding/json"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/forms"
	"github.com/adonsav/fgoapp/internal/helpers"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/adonsav/fgoapp/internal/render"
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
	handlersAppConfig *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(ac *config.AppConfig) *Repository {
	return &Repository{
		handlersAppConfig: ac,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the home page handler
func (hr *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (hr *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Registration is the registration page handler
func (hr *Repository) Registration(w http.ResponseWriter, r *http.Request) {
	var emptyRegistration models.Registration
	data := make(map[string]interface{})
	data["registration"] = emptyRegistration

	render.Template(w, r, "register.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostRegistration is the registration page post handler
func (hr *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration := models.Registration{
		UserName: r.Form.Get("user-name"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	form := forms.New(r.PostForm)
	form.Required("user-name", "email", "password")
	form.MinLength("user-name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["registration"] = registration

		render.Template(w, r, "register.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	hr.handlersAppConfig.Session.Put(r.Context(), "registration", registration)
	// It is a good practice to redirect the user to another page after a POST
	// request. This ensures that the form is not going to be submitted accidentally twice.
	http.Redirect(w, r, "/registration-summary", http.StatusSeeOther)
}

// RegistrationSummary is the registration summary page handler
func (hr *Repository) RegistrationSummary(w http.ResponseWriter, r *http.Request) {
	registration, ok := hr.handlersAppConfig.Session.Get(r.Context(), "registration").(models.Registration)
	if !ok {
		hr.handlersAppConfig.ErrorLog.Println("Can't get registration from session")
		hr.handlersAppConfig.Session.Put(r.Context(), "error", "Can't get registration from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	hr.handlersAppConfig.Session.Remove(r.Context(), "registration")

	data := make(map[string]interface{})
	data["registration"] = registration
	render.Template(w, r, "registration-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

// QuestionMark is the question mark page handler
func (hr *Repository) QuestionMark(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "question-mark.page.gohtml", &models.TemplateData{})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// QuestionMarkJSON is the question mark page JSON handler
func (hr *Repository) QuestionMarkJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "available",
	}

	output, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)

}
