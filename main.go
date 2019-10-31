package main

import (
	"encoding/json"
	"github.com/psidex/MessengerStats/messenger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

// loadJson takes a path to a json file and a Messages struct, and unmarshals the data from the file into the struct.
// Example:
//  messages := &messenger.Messages{}
//  loadJson(messagesJsonFile, messages)
func loadJson(pathToFile string, messages *messenger.Messages) {
	jsonFile, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal("Couldn't open messages file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, messages)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON data: ", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No messenger directory specified")
	}

	messagesDir := os.Args[1]
	messagesJsonFile := path.Join(messagesDir, "message_1.json")
	log.Println("Opening", messagesJsonFile)

	messages := &messenger.Messages{}
	loadJson(messagesJsonFile, messages)

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

	http.Handle("/", http.FileServer(http.Dir("static")))

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
