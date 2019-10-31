package main

import (
	"encoding/json"
	"github.com/psidex/MessengerStats/messenger"
	"log"
	"net/http"
)

var messages *messenger.Messages

// uploadMessengerFile is a HTTP endpoint that takes a JSON file and parses it into the global messages variable
func uploadMessengerFile(w http.ResponseWriter, r *http.Request) {
	// https://tutorialedge.net/golang/go-file-upload-tutorial/
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file.
	file, handler, err := r.FormFile("messengerFile")
	if err != nil {
		log.Println("Error Retrieving the File:", err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	messages = &messenger.Messages{}
	err = messenger.LoadMessengerJson(file, messages)
	if err != nil {
		// TODO: Let user know file is invalid
		log.Fatal("loadMessengerJson error:", err)
	}

	// Only redirects after data loaded into messages struct
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}

func main() {
	// Default to empty
	messages = &messenger.Messages{}

	http.HandleFunc("/api/messages/permonth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		messagesPerMonth := messenger.CountMessagesPerMonth(messages)
		_ = json.NewEncoder(w).Encode(messagesPerMonth)
	})

	http.HandleFunc("/api/messages/peruser", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		messagesPerMonth := messenger.CountMessagesPerUser(messages)
		_ = json.NewEncoder(w).Encode(messagesPerMonth)
	})

	http.HandleFunc("/api/messages/perweekday", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		messagesPerWeekday := messenger.CountMessagesPerWeekday(messages)
		_ = json.NewEncoder(w).Encode(messagesPerWeekday)
	})

	http.HandleFunc("/api/messages/title", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(messages.Title)
	})

	http.HandleFunc("/upload", uploadMessengerFile)
	http.Handle("/", http.FileServer(http.Dir("static")))
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
