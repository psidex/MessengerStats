package main

import (
	"encoding/json"
	"fmt"
	"github.com/psidex/MessengerStats/messenger"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// A struct to hold all stats from an individual Messenger conversation.
// TODO: Title can potentially have non-ASCII chars in it, how to deal with that?
type conversationData struct {
	Title          string              `json:"title"`
	MsgsPerMonth   map[int]map[int]int `json:"msgsPerMonth"`
	MsgsPerUser    map[string]int      `json:"msgsPerUser"`
	MsgsPerWeekday map[string]int      `json:"msgsPerWeekday"`
}

// A struct to hold a map of id -> conversationData.
// A mutex is used so that idLookup can be concurrently accessed (e.g. from within HTTP handlers).
type conversationDataHolder struct {
	idLookup map[string]*conversationData
	mutex    *sync.Mutex
}

// uploadMessengerFile is a HTTP handler that takes a JSON file, parses it, and inserts it into the map in `c`.
func (c *conversationDataHolder) uploadMessengerFileHandler(w http.ResponseWriter, r *http.Request) {
	// https://tutorialedge.net/golang/go-file-upload-tutorial/
	// FormFile returns the first file for the given key `myFile` it also returns the FileHeader so we can get the
	// Filename, the Header and the size of the file.
	file, handler, err := r.FormFile("messengerFile")
	if err != nil {
		log.Println("Error Retrieving the File:", err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)

	messages := &messenger.Messages{}
	err = messenger.LoadMessengerJson(file, messages)
	if err != nil {
		// TODO: Let user know file is invalid
		log.Fatal("loadMessengerJson error:", err)
	}

	// Create new ID and insert into struct lookup map.
	id := ksuid.New().String()

	// Calculating these before the Lock saves performance.
	msgsPerMonth := messenger.CountMessagesPerMonth(messages)
	msgsPerUser := messenger.CountMessagesPerUser(messages)
	msgsPerWeekday := messenger.CountMessagesPerWeekday(messages)

	c.mutex.Lock()
	c.idLookup[id] = &conversationData{
		messages.Title,
		msgsPerMonth,
		msgsPerUser,
		msgsPerWeekday,
	}
	c.mutex.Unlock()

	// Append ID to redirect URL so JS on page can get info.
	values := url.Values{}
	values.Set("id", id)

	referer := r.Header.Get("Referer")
	refererParsed, _ := url.Parse(referer)
	refererParsed.RawQuery = values.Encode()

	http.Redirect(w, r, refererParsed.String(), 302)
}

// getConversationDataApiHandler is a HTTP handler that takes an ID (param in a url) and returns the respective data from `c`.
func (c *conversationDataHolder) getConversationDataApiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idQuery, ok := r.URL.Query()["id"]

	if ok {
		id := idQuery[0]
		log.Println("Stats for ID requested:", id)

		c.mutex.Lock()
		stats, ok := c.idLookup[id]
		c.mutex.Unlock()

		if ok {
			err := json.NewEncoder(w).Encode(stats)
			if err != nil {
				log.Println("JSON Encode error:", err)
			}

		} else {
			// TODO: More consistent error reporting than just a string description?
			_, _ = fmt.Fprintf(w, "{\"error\": \"ID not found\"}")
			return
		}
	} else {
		// 404 since no ID was in the URL.
		http.NotFound(w, r)
		return
	}
}

func main() {
	convDataHolder := &conversationDataHolder{}
	convDataHolder.mutex = &sync.Mutex{}
	convDataHolder.idLookup = make(map[string]*conversationData)

	http.HandleFunc("/api/data", convDataHolder.getConversationDataApiHandler)
	http.HandleFunc("/upload", convDataHolder.uploadMessengerFileHandler)

	http.Handle("/", http.FileServer(http.Dir("static")))

	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
