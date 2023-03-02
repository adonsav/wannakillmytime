package render

import (
	"github.com/adonsav/fgoapp/internal/templates"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td templates.TemplateData
	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(request.Context(), "flash", "123")
	result := AddDefaultData(&td, request)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./internal/templates/gohtmltemplates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	renderAppConfig.TemplateCache = tc

	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var mhw myHttpWriter
	err = Template(&mhw, request, "home.page.gohtml", &templates.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&mhw, request, "non-existent.page.gohtml", &templates.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(renderAppConfig)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../templates/gohtmltemplates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/a-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx, _ = session.Load(ctx, request.Header.Get("X-Session"))
	request = request.WithContext(ctx)

	return request, nil
}
