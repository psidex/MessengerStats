package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/MessengerStats/internal/api"
	"log"
	"net/http"
)

func SendView(htmlFile string) func(w http.ResponseWriter, r *http.Request) {
	// TODO: Doing it this way means the file is opened and read on every request, find way to load into mem.
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./views/"+htmlFile+".html")
	}
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", SendView("index"))
	r.HandleFunc("/stats", SendView("stats"))

	statsApi := api.NewStatsApi()
	r.HandleFunc("/api/stats", statsApi.StatsHandler).Methods("GET")
	r.HandleFunc("/api/upload", statsApi.UploadHandler).Methods("POST")

	log.Print("Serving at http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
