// Web package provides:
//   - handlers to handle http requests;
//   - http server interface to start up it and stop it.
//
// Use [ConnectAll] to access final handler.
// Handlers logic handlers are in the sub-packages.
//
// Before usage, call [github.com/TrueHopolok/VladOS/modules/web/webtmls.PrepareTemplates] for html templates to load.
// Then you can use [Start] and/or [Stop] to control the server.
package web

import (
	"net/http"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

var server *http.Server

// Connects all connectors into 1 new [net/http.ServeMux] to serve.
//   - [LoggerMiddleware],
//   - [ConnectEveryone],
//   - [ConnectAuthorized],
//   - [ConnectUnauthorized],
//   - [ConnectFileHandlers].
func NewWebHandler() http.Handler {
	mux := http.NewServeMux()
	ConnectEveryone(mux)
	ConnectAuthorized(mux)
	ConnectUnauthorized(mux)
	ConnectFileHandlers(mux)
	return LoggerMiddleware(mux)
}

// Start the http server as a separate goroutinue.
// Will close the existing server if it was opened by this package.
//
// Server options are:
//
//	&http.Server{
//		Addr:    ":8080",
//		Handler: NewWebHandler(),
//	}
//
// Returns error if happens on initialization.
// Otherwise uses provided channel to report about the error while executing the server.
func Start(serverErrorChan chan error) error {
	if server != nil {
		server.Close()
	}
	server = &http.Server{
		Addr:    ":8080",
		Handler: NewWebHandler(),
	}
	go func() {
		serverErrorChan <- server.ListenAndServe()
	}()
	select {
	case err := <-serverErrorChan:
		return err
	default:
		return nil
	}
}

// Stops existing server from executing. If no server was opened, will do nothing.
// Returns an error from [net/http.Server.Close].
func Stop() error {
	if server == nil {
		return nil
	}
	return server.Close()
}
