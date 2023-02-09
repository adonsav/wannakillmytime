package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"GETHome", "/", "GET", []postData{}, http.StatusOK},
	{"GETAbout", "/about", "GET", []postData{}, http.StatusOK},
	{"GETRegister", "/register", "GET", []postData{}, http.StatusOK},
	{"GETRegSummary", "/registration-summary", "GET", []postData{}, http.StatusOK},
	{"GETQueMark", "/question-mark", "GET", []postData{}, http.StatusOK},
	{"POSTRegister", "/register", "POST", []postData{
		{key: "user-name", value: "test"},
		{key: "email", value: "test@test.com"},
	}, http.StatusOK},
	{"POSTQueMarkJSON", "/question-mark-json", "POST", []postData{
		{key: "user-name", value: "test"},
		{key: "email", value: "test@test.com"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		if test.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
