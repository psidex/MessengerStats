package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/MessengerStats/internal/handlers"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	indexBytes, err := ioutil.ReadFile("./views/index.html")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", handlers.SendBytes(indexBytes)).Methods("GET")
	r.HandleFunc("/stats", handlers.Redirect("/")).Methods("GET")
	r.HandleFunc("/api/ws", handlers.WebSocketApi).Methods("GET")

	log.Print("Serving at http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
