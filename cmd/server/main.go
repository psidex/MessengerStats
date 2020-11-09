package main

import (
	"github.com/gorilla/mux"
	"github.com/psidex/MessengerStats/internal/handlers"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stats", handlers.Redirect("/")).Methods("GET")
	r.HandleFunc("/api/ws", handlers.WebSocketApi).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	log.Print("Serving http on all available interfaces @ port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
