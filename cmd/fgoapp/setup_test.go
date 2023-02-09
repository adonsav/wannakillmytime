package main

import (
	"net/http"
	"os"
	"testing"
)

type testingHandler struct{}

func (th *testingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

