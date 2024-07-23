// Package http includes handlers and utilties.
package http

import "net/http"

// JSONVersionHandler returns a simple JSON response from a version string.
// The returned JSON is in the form of `{"version":"v0.0.1"}`.
func JSONVersionHandler(version string) http.HandlerFunc {
	bodyBytes := []byte(`{"version":"` + version + `"}`)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bodyBytes)
	}
}
