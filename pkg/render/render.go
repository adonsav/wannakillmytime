package render

import (
	"github.com/adonsav/fgoapp/pkg/config"
	"github.com/adonsav/fgoapp/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var appConfig *config.AppConfig

// NewTemplates sets the configuration for the templates
func NewTemplates(ac *config.AppConfig) {
	appConfig = ac
}

// Template renders a template
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template
	if appConfig.UseCache {
		// get the template cache from application configuration
		templateCache = appConfig.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	templ, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	td = AddDefaultData(td)

	err := templ.Execute(w, td)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files name *.page.gohtml from ./templates
	pages, err := filepath.Glob("./templates/*.page.gohtml")
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

		layouts, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}

	return myCache, nil
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}
