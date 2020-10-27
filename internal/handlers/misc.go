package handlers

import "net/http"

// SendBytes takes the given bytes and sends them on every request.
func SendBytes(bytes []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(bytes)
	}
}

// Redirect redirects every request to the given path.
func Redirect(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusMovedPermanently)
	}
}
