package web

import (
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webindex"
)

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Everyone]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectEveryone(mux *http.ServeMux) {
	mux.HandleFunc("GET /", vos.AuthMiddlewareFunc(webindex.TodoHandler, vos.Everyone))
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
		http.ServeFile(w, r, "./static/favicon.png")
	})
}
