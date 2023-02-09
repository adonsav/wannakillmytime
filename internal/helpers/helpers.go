package helpers

import (
	"fmt"
	"github.com/adonsav/fgoapp/internal/config"
	"net/http"
	"runtime/debug"
)

var helpersAppConfig *config.AppConfig

// NewHelpers sets up application configuration for helpers
func NewHelpers(ac *config.AppConfig) {
	helpersAppConfig = ac
}

func ClientError(w http.ResponseWriter, status int) {
	helpersAppConfig.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	helpersAppConfig.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
