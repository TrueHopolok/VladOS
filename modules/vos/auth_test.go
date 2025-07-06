package vos_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

// Made to do nothing for the request.
// Used to test authefication middleware.
type emptyHandler struct{}

func (emptyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { /*do nothing*/ }

// Creates an server for testing [vos.AuthMiddleware].
// Only executes middleware and inside handler does nothing.
func initCommunication() (*httptest.Server, *http.Client) {
	var handler emptyHandler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /everyone/", vos.AuthMiddleware(handler, vos.Everyone))
	mux.HandleFunc("GET /authorized/", vos.AuthMiddleware(handler, vos.Authorized))
	mux.HandleFunc("GET /unauthorized/", vos.AuthMiddleware(handler, vos.Unauthorized))
	server := httptest.NewServer(mux)
	return server, server.Client()
}

func TestEveryone(t *testing.T) {
	// server, client := initCommunication()
}

func TestAuthorized(t *testing.T) {
	// server, client := initCommunication()
}

func TestUnauthorizes(t *testing.T) {
	// server, client := initCommunication()
}
