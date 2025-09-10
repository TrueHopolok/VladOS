package vos_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/mlog"
	"github.com/TrueHopolok/VladOS/modules/vos"
)

const pathToRoot = "../../"

// Made to do nothing for the request.
// Used to test authefication middleware.
type emptyHandler struct{}

func (emptyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { /*do nothing*/ }

// Creates an server for testing [vos.AuthMiddleware].
// Only executes middleware and inside handler does nothing.
func initServer() *httptest.Server {
	var handler emptyHandler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /everyone/", vos.AuthMiddleware(handler, vos.Everyone))
	mux.HandleFunc("GET /authorized/", vos.AuthMiddleware(handler, vos.Authorized))
	mux.HandleFunc("GET /unauthorized/", vos.AuthMiddleware(handler, vos.Unauthorized))
	return httptest.NewServer(mux)
}

// Adds valid auth cookie to the client for future requests.
func getNewAuthCookie(t *testing.T, server *httptest.Server) *http.Cookie {
	url, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("unexpected url unparsing error: %s", err)
	}
	if url == nil {
		t.Fatalf("received url is nil")
	}
	ses := vos.NewSession(0, "test")
	jwt, err := ses.NewJWT()
	if err != nil {
		t.Fatalf("unexpected jwt generation error: %s", err)
	}
	return &http.Cookie{
		Name:     vos.AuthCookieName,
		Value:    jwt,
		Path:     "/",
		MaxAge:   int(vos.AuthExpires.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}
}

func testTemplate(t *testing.T, path string, firstCode, secondCode int) {
	mlog.InitTesting(t, pathToRoot) // used since auth logs its actions
	server := initServer()
	path = server.URL + path
	client := server.Client()

	req1, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("unexpected 1st request init error: %s", err)
	}
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("http request failed; err: %s", err)
	}
	if resp1.StatusCode != firstCode {
		msg, _ := io.ReadAll(resp1.Body)
		t.Fatalf("unexpected 2nd status code\nwant: %d\ngot: %d\nmsg: %s", secondCode, resp1.StatusCode, string(msg))
	}
	resp1.Body.Close()

	req2, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatalf("unexpected 2nd request init error: %s", err)
	}
	req2.AddCookie(getNewAuthCookie(t, server))
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("http request failed; err: %s", err)
	}
	if resp2.StatusCode != secondCode {
		msg, _ := io.ReadAll(resp2.Body)
		t.Fatalf("unexpected 2nd status code\nwant: %d\ngot: %d\nmsg: %s", secondCode, resp2.StatusCode, string(msg))
	}
	resp2.Body.Close()
}

func TestEveryone(t *testing.T) {
	testTemplate(t, "/everyone/", http.StatusOK, http.StatusOK)
}

func TestAuthorized(t *testing.T) {
	testTemplate(t, "/authorized/", http.StatusUnauthorized, http.StatusOK)
}

func TestUnauthorized(t *testing.T) {
	testTemplate(t, "/unauthorized/", http.StatusOK, http.StatusBadRequest)
}
