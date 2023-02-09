package render

import (
	"errors"
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var renderAppConfig *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates sets the configuration for the templates
func NewTemplates(ac *config.AppConfig) {
	renderAppConfig = ac
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template
	if renderAppConfig.UseCache {
		// get the template cache from application configuration
		templateCache = renderAppConfig.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	templ, ok := templateCache[tmpl]
	if !ok {
		return errors.New("could not get template from template cache")

	}

	td = AddDefaultData(td, r)

	err := templ.Execute(w, td)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files name *.page.gohtml from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.gohtml
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
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

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = renderAppConfig.Session.PopString(r.Context(), "flash")
	td.Error = renderAppConfig.Session.PopString(r.Context(), "error")
	td.Warning = renderAppConfig.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}
