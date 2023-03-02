package handlers

import (
	"context"
	"encoding/json"
	"github.com/adonsav/fgoapp/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"register", "/register", "GET", http.StatusOK},
	{"regSummary", "/registration-summary", "GET", http.StatusOK},
	{"queMark", "/question-mark", "GET", http.StatusOK},
	{"non-existent", "/non/existent", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"activeRegs", "/admin/registrations-active", "GET", http.StatusOK},
	{"allRegs", "/admin/registrations-all", "GET", http.StatusOK},
	{"showReg", "/admin/registrations/all/1", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		resp, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_RegistrationSummary(t *testing.T) {
	registration := models.Registration{
		ID:       -1,
		UserName: "TestUserName",
		Email:    "testEmail@email.com",
		Password: "TestPassword",
		Active:   true,
	}

	request, _ := http.NewRequest("GET", "/registration-summary", nil)
	ctx := getCtx(request)
	request = request.WithContext(ctx)

	// httptest.NewRecorder() simulates the request-response lifecycle and satisfies
	// the requirements for http.ResponseWriter
	responseRecorder := httptest.NewRecorder()
	session.Put(ctx, "registration", registration)

	// http.HandlerFunc allows us to call a handler function directly as a normal function
	handler := http.HandlerFunc(Repo.RegistrationSummary)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("RegistrationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// Case where registration is missing from session
	request, _ = http.NewRequest("GET", "/registration-summary", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)
	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("RegistrationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostRegistration(t *testing.T) {
	requestBody := url.Values{}
	requestBody.Add("user-name", "testUserName")
	requestBody.Add("email", "testEmail@email.com")
	requestBody.Add("password", "testPassword")

	request, _ := http.NewRequest("POST", "/register", strings.NewReader(requestBody.Encode()))
	ctx := getCtx(request)
	request = request.WithContext(ctx)

	// It is not required to set the header, but it is considered a good practice.
	// This says to the web server that the request it's going to get is a form POST.
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostRegistration)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostRegistration handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case where there is a missing post body
	request, _ = http.NewRequest("POST", "/register", nil)
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostRegistration handler returned wrong response code for missing post body: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case where the form has invalid input data
	requestBody.Set("user-name", "X")
	requestBody.Add("email", "testEmail@email.com")
	requestBody.Add("password", "testPassword")

	request, _ = http.NewRequest("POST", "/register", strings.NewReader(requestBody.Encode()))
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("PostRegistration handler returned wrong response code for invalid form data: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case where the insertion of a registration into database fails
	requestBody.Set("user-name", "failRegistration")
	requestBody.Add("email", "testEmail@email.com")
	requestBody.Add("password", "testPassword")

	request, _ = http.NewRequest("POST", "/register", strings.NewReader(requestBody.Encode()))
	ctx = getCtx(request)
	request = request.WithContext(ctx)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostRegistration handler failed on intended registration insertion failure: got %d, wanted %d",
			responseRecorder.Code,
			http.StatusSeeOther)
	}
}

func TestRepository_QuestionMarkJSON(t *testing.T) {
	requestBody := url.Values{}
	requestBody.Add("so-what", "test")
	requestBody.Add("see-above", "test")

	request, _ := http.NewRequest("POST", "/question-mark-json", strings.NewReader(requestBody.Encode()))

	ctx := getCtx(request)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.QuestionMarkJSON)
	handler.ServeHTTP(responseRecorder, request)

	var jsonResp jsonResponse
	err := json.Unmarshal([]byte(responseRecorder.Body.String()), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	// Case where there is an error parsing the form
	request, _ = http.NewRequest("POST", "/question-mark-json", nil)

	ctx = getCtx(request)
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.QuestionMarkJSON)
	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal([]byte(responseRecorder.Body.String()), &jsonResp)
	if err != nil {
		t.Error("failed to parse json!")
	}

	// Since we specified a nil body, we expect an internal server error
	if jsonResp.Message != "internal server error" {
		t.Error("Got proper response when request body was empty")
	}
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"validCredentials",
		"validEmail@test.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"nonExistingUser",
		"nonExistingUser@test.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalidData",
		"email",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestRepository_Login(t *testing.T) {
	for _, test := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", test.email)
		postedData.Add("password", "password")

		request, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(request)
		request = request.WithContext(ctx)

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostLogin)
		handler.ServeHTTP(responseRecorder, request)

		if responseRecorder.Code != test.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d",
				test.name, test.expectedStatusCode, responseRecorder.Code)
		}

		if test.expectedLocation != "" {
			actualLocation, _ := responseRecorder.Result().Location()
			if actualLocation.String() != test.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got %s",
					test.name, test.expectedLocation, actualLocation.String())
			}
		}

		if test.expectedHTML != "" {
			html := responseRecorder.Body.String()
			if !strings.Contains(html, test.expectedHTML) {
				t.Errorf("failed %s: expected to find %s, but got %s",
					test.name, test.expectedHTML, html)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
