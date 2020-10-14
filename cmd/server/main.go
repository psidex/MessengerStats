package main

import (
	"github.com/psidex/MessengerStats/internal/api"
	"net/http"
)

func main() {
	statsApi := api.NewConversationStatsApi()
	http.HandleFunc("/api/stats", statsApi.ConversationStatsHandler)
	http.HandleFunc("/upload", statsApi.FileUploadHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}