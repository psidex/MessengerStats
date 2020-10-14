package api

import (
	"encoding/json"
	"fmt"
	"github.com/psidex/MessengerStats/internal/messenger"
	"github.com/psidex/MessengerStats/internal/stats"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// conversationStats holds all the statistics for an individual Messenger conversation.
type conversationStats struct {
	Title              string              `json:"conversation_title"`
	MessagesPerMonth   map[int]map[int]int `json:"messages_per_month"`
	MessagesPerUser    map[string]int      `json:"messages_per_user"`
	MessagesPerWeekday map[string]int      `json:"messages_per_weekday"`
}

// ConversationStatsApi contains the data and functions for generating conversation statistics.
type ConversationStatsApi struct {
	savedStats map[string]*conversationStats // [id]data
	mu         *sync.Mutex
}

func NewConversationStatsApi() *ConversationStatsApi {
	c := &ConversationStatsApi{}
	c.mu = &sync.Mutex{}
	c.savedStats = make(map[string]*conversationStats)
	return c
}

// FileUploadHandler is a HTTP handler that takes a Messenger JSON file, parses it, and saves the stats in memory.
// Expects a POST request.
func (c *ConversationStatsApi) FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	// https://tutorialedge.net/golang/go-file-upload-tutorial/
	file, fileHeader, err := r.FormFile("messengerFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "error: %s}", err)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	log.Printf("File Size: %+v\n", fileHeader.Size)

	conversation, err := messenger.NewConversation(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprintf(w, "error: %s}", err)
		return
	}

	id := ksuid.New().String()
	startTime := time.Now()

	mpuCounter := stats.NewMessagesPerUserCounter()
	mpdCounter := stats.NewMessagesPerWeekdayCounter()
	mpmCounter := stats.NewMessagesPerMonthCounter()

	for _, message := range conversation.Messages {
		mpuCounter.Update(message)
		mpdCounter.Update(message)
		mpmCounter.Update(message)
	}

	log.Printf("Statistics calculations took %s", time.Since(startTime))

	c.mu.Lock()
	c.savedStats[id] = &conversationStats{
		conversation.Title,
		mpmCounter.MessagesPerYearMonth,
		mpuCounter.MessagesPerUser,
		mpdCounter.MessagesPerDay,
	}
	c.mu.Unlock()

	// Send the user from whence they came (with the id that they need).
	values := url.Values{}
	values.Set("id", id)
	referer := r.Header.Get("Referer")
	refererParsed, _ := url.Parse(referer)
	refererParsed.RawQuery = values.Encode()

	http.Redirect(w, r, refererParsed.String(), 302)
}

// ConversationStatsHandler is a HTTP handler that takes an ID (param in a url) and returns the respective data from `c`.
func (c *ConversationStatsApi) ConversationStatsHandler(w http.ResponseWriter, r *http.Request) {
	idQuery, ok := r.URL.Query()["id"]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	id := idQuery[0]
	log.Println("Stats for ID requested:", id)

	c.mu.Lock()
	savedStats, ok := c.savedStats[id]
	c.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "{\"error\": \"ID not found\"}")
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(savedStats)
	if err != nil {
		log.Println("JSON Encode error:", err)
	}
}
