// Web package provides:
//   - handlers to handle http requests;
//   - http server interface to start up it and stop it.
//
// Use [ConnectAll] to access final handler.
// Handlers logic handlers are in the sub-packages.
//
// Use [Start] and/or [Stop] to control the server.
package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	"github.com/TrueHopolok/VladOS/modules/vos"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

var server *http.Server

// Provides a placeholder for a handler function to be used in the webpage as TODO reminder.
func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("met=%s  url=%v\nHere is nothing...", r.Method, r.URL)))
}

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Everyone]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectEveryone(mux *http.ServeMux) {
	mux.HandleFunc("GET /", vos.AuthMiddlewareFunc(TodoHandler, vos.Everyone))
}

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Authorized]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectAuthorized(mux *http.ServeMux) {
	// connect all handlers
}

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Unauthorized]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectUnauthorized(mux *http.ServeMux) {
	// connect all handlers
}

// Connects to [net/http.ServeMux] handler functions to server static files
// on the [github.com/TrueHopolok/VladOS/modules/cfg.Config.WebStaticDir] path.
func ConnectFileHandlers(mux *http.ServeMux) {
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.Get().WebStaticDir))))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static.favicon")
	})
}

// Connects all connector into 1 new [net/http.ServeMux] to serve.
//
// For connectors info see [ConnectEveryone], [ConnectAuthorized], [ConnectUnauthorized] and [ConnectFileHandlers].
func NewWebHandler() http.Handler {
	mux := http.NewServeMux()
	ConnectEveryone(mux)
	ConnectAuthorized(mux)
	ConnectUnauthorized(mux)
	ConnectFileHandlers(mux)
	return LoggerMiddleware(mux)
}

// Provides small http middleware for logs purposes using [log/slog] package.
func LoggerMiddleware(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("http exe", "mtd", r.Method, "url", r.URL, "status", "START")
		defer slog.Debug(r.Method, "url", r.URL, "status", "FINISH")
		defer func() {
			if x := recover(); x != nil {
				slog.Error(r.Method, "url", r.URL, "status", "FAILED", "panic", x)
				http.Error(w, "http handler paniced", http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	}
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
