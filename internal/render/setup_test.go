package render

import (
	"encoding/gob"
	"github.com/adonsav/fgoapp/internal/config"
	"github.com/adonsav/fgoapp/internal/models"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var renderTestAppConfig config.AppConfig

func TestMain(m *testing.M) {
	renderTestAppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	renderTestAppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	gob.Register(models.Registration{})
	renderTestAppConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	renderTestAppConfig.Session = session

	renderAppConfig = &renderTestAppConfig

	os.Exit(m.Run())
}

type myHttpWriter struct{}

func (mhw *myHttpWriter) Header() http.Header {
	return http.Header{}
}

func (mhw *myHttpWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (mhw *myHttpWriter) WriteHeader(statusCode int) {}
