package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

// Provides a placeholder for a handler function to be used in the webpage as TODO reminder.
func TodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("met=%s  url=%v\nHere is nothing...", r.Method, r.URL)))
}

// Provides small http middleware for logs purposes using [log/slog] package.
func LoggerMiddleware(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "status", "START")
		defer slog.Debug("http req", "mtd", r.Method, "url", r.URL, "status", "FINISH")
		defer func() {
			if x := recover(); x != nil {
				slog.Error("http req", "mtd", r.Method, "url", r.URL, "status", "FAILED", "panic", x)
				http.Error(w, "http handler paniced", http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	}
}

// Provides an check for being an admin before accessing the handler.
// Must be used after [vos.AuthMiddleware] was used, otherwise will block all trafic.
func AdminMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ses, auth := vos.GetSession(r)
		if !auth {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "stopped", "not autheficated")
			http.Error(w, "access denied - not autheficated", http.StatusBadRequest)
			return
		} else if !ses.Admin {
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "stopped", "not admin")
			http.Error(w, "access denied - not admin", http.StatusBadRequest)
			return
		}
		handler(w, r)
	}
}
