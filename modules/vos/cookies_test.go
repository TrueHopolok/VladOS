package vos_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

func TestGetCookie(t *testing.T) {
	ses := vos.NewSession(0, "test", true)
	jwt, err := ses.NewJWT()
	if err != nil {
		t.Fatalf("generating new jwt failed: %s", err)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:     vos.AuthCookieName,
		Value:    jwt,
		Path:     "/",
		MaxAge:   int(vos.AuthExpires.Seconds()),
		HttpOnly: true,
		Secure:   true,
	})
	token, err := vos.GetAuthCookie(r)
	if err != nil {
		t.Fatalf("unexpected error while getting cookie: %s", err)
	}
	if jwt != token {
		t.Fatalf("corrupted jwt token:\nwant: %s\ngot: %s\n", jwt, token)
	}
}

func TestSetCookie(t *testing.T) {
	ses := vos.NewSession(0, "test", true)
	jwt, err := ses.NewJWT()
	if err != nil {
		t.Fatalf("generating new jwt failed: %s", err)
	}

	w := httptest.NewRecorder()
	vos.SetAuthCookie(w, jwt)
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("invalid cookies amount")
	}
	if cookies[0].Name != vos.AuthCookieName {
		t.Fatalf("invalid cookie name")
	}
	token := cookies[0].Value
	got, valid, err := vos.ValidateJWT(token)
	if err != nil {
		t.Fatalf("validating jwt error: %s", err)
	}
	if !valid {
		t.Fatalf("received session is invalid, when should be the opposite")
	}
	if !equalSessions(ses, got) {
		t.Fatalf("recieved session is corrupted:\nwant: %v\ngot: %v\n", ses, got)
	}
}

func TestDeleteCookie(t *testing.T) {
	w := httptest.NewRecorder()
	vos.DeleteAuthCookie(w)
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("invalid cookies amount")
	}
	if cookies[0].Name != vos.AuthCookieName {
		t.Fatalf("invalid cookie name")
	}
	if cookies[0].MaxAge > 0 {
		t.Fatalf("cookie's maxage is invalid: won't delete the cookie")
	}
}
