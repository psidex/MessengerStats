package handlers

import "net/http"

// Redirect redirects every request to the given path.
func Redirect(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
}
