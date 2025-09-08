package web

import (
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/page_index"
	"github.com/TrueHopolok/VladOS/modules/web/page_leaderboard"
	"github.com/TrueHopolok/VladOS/modules/web/page_login"
	"github.com/TrueHopolok/VladOS/modules/web/page_suggestions"
)

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Everyone]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectEveryone(mux *http.ServeMux) {
	mux.HandleFunc("GET /", vos.AuthMiddlewareFunc(page_index.Handle, vos.Everyone))
	mux.HandleFunc("GET /leaderboard", vos.AuthMiddlewareFunc(page_leaderboard.Handle, vos.Everyone))
	mux.HandleFunc("GET /suggestions", vos.AuthMiddlewareFunc(page_suggestions.Handle, vos.Everyone))
}

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Authorized]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectAuthorized(mux *http.ServeMux) {
	// add pages with vox.Authorized
}

// Connects to [net/http.ServeMux] handler functions with
// permission flag [github.com/TrueHopolok/VladOS/modules/vos.Unauthorized]
// for the function [github.com/TrueHopolok/VladOS/modules/vos.AuthMiddleware].
func ConnectUnauthorized(mux *http.ServeMux) {
	mux.HandleFunc("GET /login", vos.AuthMiddlewareFunc(page_login.Handle, vos.Unauthorized))
}

// Connects to [net/http.ServeMux] handler functions to server static files
// on the [github.com/TrueHopolok/VladOS/modules/cfg.Config.WebStaticDir] path.
func ConnectFileHandlers(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir(cfg.Get().WebStaticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.Handle("GET /login/static/", http.StripPrefix("/login/static/", fs))
	mux.Handle("GET /leaderboard/static/", http.StripPrefix("/leaderboard/static/", fs))
	mux.Handle("GET /suggestions/static/", http.StripPrefix("/suggestions/static/", fs))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.png")
	})
}
