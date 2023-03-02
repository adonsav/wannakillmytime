package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/driver"
	"github.com/adonsav/fgoapp/internal/forms"
	"github.com/adonsav/fgoapp/internal/helpers"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/adonsav/fgoapp/internal/render"
	"github.com/adonsav/fgoapp/internal/repository"
	"github.com/adonsav/fgoapp/internal/templates"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Repository is the repository type
type Repository struct {
	handlersAppConfig *config.AppConfig
	DB                repository.DatabaseRepo
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(ac *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		handlersAppConfig: ac,
		DB:                repository.NewPostgresDBRepo(db.SQL, ac),
	}
}

// NewTestRepo creates a new testing repository
func NewTestRepo(ac *config.AppConfig) *Repository {
	return &Repository{
		handlersAppConfig: ac,
		DB:                repository.NewTestingDBRepo(ac),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the home page handler
func (hr *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &templates.TemplateData{})
}

// About is the about page handler
func (hr *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.gohtml", &templates.TemplateData{})
}

// Registration is the registration page handler
func (hr *Repository) Registration(w http.ResponseWriter, r *http.Request) {
	var emptyRegistration models.Registration
	data := make(map[string]interface{})
	data["registration"] = emptyRegistration

	render.Template(w, r, "registration.page.gohtml", &templates.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostRegistration is the registration page post handler
func (hr *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	// it is considered a good practice to always parse the form
	err := r.ParseForm()
	if err != nil {
		hr.handlersAppConfig.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		render.Template(w, r, "registration.page.gohtml", &templates.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = hr.DB.InsertRegistration(registration)
	if err != nil {
		hr.handlersAppConfig.Session.Put(r.Context(), "error", "can't insert registration into database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	msg := models.EmailData{
		To:       registration.Email,
		From:     "me@here.co",
		Subject:  "Registration confirmation",
		Content:  fmt.Sprintf("<strong>%s</strong> welcome to this non-sense boring site!", registration.UserName),
		Template: "basic.html",
	}
	hr.handlersAppConfig.EmailChan <- msg

	hr.handlersAppConfig.Session.Put(r.Context(), "registration", registration)
	// It is a good practice to redirect the user to another page after a POST
	// request. This ensures that the form is not going to be submitted accidentally twice.
	http.Redirect(w, r, "/registration-summary", http.StatusSeeOther)
}

// RegistrationSummary is the registration summary page handler
func (hr *Repository) RegistrationSummary(w http.ResponseWriter, r *http.Request) {
	registration, ok := hr.handlersAppConfig.Session.Get(r.Context(), "registration").(models.Registration)
	if !ok {
		hr.handlersAppConfig.Session.Put(r.Context(), "error", "can't get registration from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	hr.handlersAppConfig.Session.Remove(r.Context(), "registration")

	data := make(map[string]interface{})
	data["registration"] = registration
	render.Template(w, r, "registration-summary.page.gohtml", &templates.TemplateData{
		Data: data,
	})
}

// QuestionMark is the question mark page handler
func (hr *Repository) QuestionMark(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "question-mark.page.gohtml", &templates.TemplateData{})
}

type jsonResponse struct {
	SoWhat   string `json:"so_what"`
	SeeAbove string `json:"see_above"`
	Message  string `json:"message"`
}

// QuestionMarkJSON is the question mark page JSON handler that sends a JSON response
func (hr *Repository) QuestionMarkJSON(w http.ResponseWriter, r *http.Request) {
	// Parse the form is considered a good practice, but it is not actually needed for the
	// application to work. But it is mandatory for testing purposes.
	err := r.ParseForm()
	if err != nil {
		// Can't parse form, so return appropriate json.
		resp := jsonResponse{
			SoWhat:   "",
			SeeAbove: "",
			Message:  "internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		SoWhat:   r.Form.Get("so-what"),
		SeeAbove: r.Form.Get("see-above"),
		Message:  "OK",
	}

	// I removed the error check, since we handle all aspects of the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Login is the show login page handler
func (hr *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &templates.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin is the login page post handler
func (hr *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// It is considered a good practice to renew the token after a POST request for security reasons
	_ = hr.handlersAppConfig.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	userCredentials := models.User{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.gohtml", &templates.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := hr.DB.Authenticate(userCredentials.Email, userCredentials.Password)
	if err != nil {
		log.Println(err)
		hr.handlersAppConfig.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	hr.handlersAppConfig.Session.Put(r.Context(), "user_id", id)
	hr.handlersAppConfig.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout is the logout page handler
func (hr *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = hr.handlersAppConfig.Session.Destroy(r.Context())
	_ = hr.handlersAppConfig.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard is the admin dashboard page handler
func (hr *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.gohtml", &templates.TemplateData{})
}

// AdminActiveRegistrations is the admin active registrations page handler
func (hr *Repository) AdminActiveRegistrations(w http.ResponseWriter, r *http.Request) {
	registrations, err := hr.DB.ActiveRegistrations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["registrations"] = registrations

	render.Template(w, r, "admin-active-registrations.page.gohtml", &templates.TemplateData{
		Data: data,
	})
}

// AdminAllRegistrations is the admin all registrations page handler
func (hr *Repository) AdminAllRegistrations(w http.ResponseWriter, r *http.Request) {
	registrations, err := hr.DB.AllRegistrations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["registrations"] = registrations

	render.Template(w, r, "admin-all-registrations.page.gohtml", &templates.TemplateData{
		Data: data,
	})
}

// AdminShowRegistration is the admin show registration page handler
func (hr *Repository) AdminShowRegistration(w http.ResponseWriter, r *http.Request) {
	splittedURL := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(splittedURL[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := splittedURL[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	registration, err := hr.DB.GetRegistrationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["registration"] = registration

	render.Template(w, r, "admin-show-registration.page.gohtml", &templates.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

// AdminPostShowRegistration is the admin show registration page post handler
func (hr *Repository) AdminPostShowRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	splittedURL := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(splittedURL[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := splittedURL[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	registration, err := hr.DB.GetRegistrationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	registration.UserName = r.Form.Get("user-name")
	registration.Email = r.Form.Get("email")

	err = hr.DB.UpdateRegistration(registration)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	hr.handlersAppConfig.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/registrations-%s", src), http.StatusSeeOther)
	render.Template(w, r, "admin-show-registration.page.gohtml", &templates.TemplateData{

		Form: forms.New(nil),
	})
}

// AdminDeactivateRegistration is the admin deactivate registration page post handler
func (hr *Repository) AdminDeactivateRegistration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	registration, err := hr.DB.GetRegistrationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if !registration.Active {
		hr.handlersAppConfig.Session.Put(r.Context(), "warning", "Registration already deactivated.")
	} else {
		_ = hr.DB.DeactivateRegistration(id)
		hr.handlersAppConfig.Session.Put(r.Context(), "flash", "Registration deactivated!")
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/registrations-%s", src), http.StatusSeeOther)
}
