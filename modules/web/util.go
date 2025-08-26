package web

import (
	"fmt"
	"log/slog"
	"net/http"
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
