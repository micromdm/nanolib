// Package http includes handlers and utilties.
package http

import (
	"crypto/subtle"
	"net/http"
)

// NewJSONVersionHandler responds with JSON from version.
// The JSON is in the form of `{"version":"v0.0.1"}`.
// No validation is done on version; it should conform to a JSON string.
func NewJSONVersionHandler(version string) http.HandlerFunc {
	bodyBytes := []byte(`{"version":"` + version + `"}`)
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bodyBytes)
	}
}

// NewSimpleBasicAuthHandler is an HTTP Basic authentication middleware.
// The returned handler verifies the the correct authentication against
// username and password before handing off to next, or reponds with
// an HTTP 401 setting the WWW-Authenticate header using realm.
func NewSimpleBasicAuthHandler(next http.Handler, username, password, realm string) http.HandlerFunc {
	// cache 1-time data
	ubc := []byte(username)
	pbc := []byte(password)
	rc := `Basic realm="` + realm + `"`
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(u), ubc) != 1 ||
			subtle.ConstantTimeCompare([]byte(p), pbc) != 1 {
			w.Header().Set("Www-Authenticate", rc)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
