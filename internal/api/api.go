package api

import (
	"encoding/json"
	"fmt"
	"github.com/psidex/MessengerStats/internal/messenger"
	"github.com/psidex/MessengerStats/internal/stats"
	"github.com/segmentio/ksuid"
	"log"
	"mime/multipart"
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
	savedStats map[string]*conversationStats // Map of generated id : struct pointer
	mu         *sync.Mutex                   // Controls access to savedStats
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
	if r.Method != "POST" {
		http.Error(w, "must send a POST request to this endpoint", http.StatusBadRequest)
		return
	}

	const _500MB = 1024 * 1024 * 500 // Use up to 500MB of RAM, rest goes on disk.
	if err := r.ParseMultipartForm(_500MB); nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mpuCounter := stats.NewMessagesPerUserCounter()
	mpdCounter := stats.NewMessagesPerWeekdayCounter()
	mpmCounter := stats.NewMessagesPerMonthCounter()

	title := ""
	id := ksuid.New().String()
	startTime := time.Now()

	headers, ok := r.MultipartForm.File["messenger_files"]
	if !ok {
		http.Error(w, "invalid request, could not find messenger_files in form", http.StatusInternalServerError)
		return
	}

	for _, header := range headers {
		var (
			file multipart.File
			err  error
		)

		if file, err = header.Open(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conversation, err := messenger.NewConversation(file)
		_ = file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		title = conversation.Title

		for _, message := range conversation.Messages {
			mpuCounter.Update(message)
			mpdCounter.Update(message)
			mpmCounter.Update(message)
		}
	}

	log.Printf("File parse and calculations took %s", time.Since(startTime))

	c.mu.Lock()
	c.savedStats[id] = &conversationStats{
		title,
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
		_, _ = fmt.Fprint(w, "{\"error\": \"ID not found\"}")
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(savedStats)
	if err != nil {
		log.Println("ConversationStatsHandler JSON Encode error:", err)
	}
}
