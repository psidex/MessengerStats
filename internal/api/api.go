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
	"sync"
	"time"
)

// apiResponse holds all the data for a conversation.
type apiResponse struct {
	Title              string                           `json:"conversation_title"`
	MessagesPerMonth   stats.MessagesPerMonthJsObject   `json:"messages_per_month"`
	MessagesPerUser    []stats.MessagesPerUserJsObject  `json:"messages_per_user"`
	MessagesPerWeekday stats.MessagesPerWeekdayJsObject `json:"messages_per_weekday"`
}

// ConversationStatsApi contains the data and functions for generating conversation statistics.
type ConversationStatsApi struct {
	apiResponseCache map[string]*apiResponse // unique id : cache
	mu               *sync.Mutex             // Controls access to apiResponseCache
}

func NewConversationStatsApi() *ConversationStatsApi {
	c := &ConversationStatsApi{}
	c.mu = &sync.Mutex{}
	c.apiResponseCache = make(map[string]*apiResponse)
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

	headers, ok := r.MultipartForm.File["messenger_files"]
	if !ok {
		http.Error(w, "invalid request, could not find messenger_files in form", http.StatusInternalServerError)
		return
	}

	mpuCounter := stats.NewMessagesPerUserCounter()
	mpwdCounter := stats.NewMessagesPerWeekdayCounter()
	mpmCounter := stats.NewMessagesPerMonthCounter()

	title := ""
	id := ksuid.New().String()
	startTime := time.Now()

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
			mpwdCounter.Update(message)
			mpmCounter.Update(message)
		}
	}

	c.mu.Lock()
	c.apiResponseCache[id] = &apiResponse{
		title,
		mpmCounter.GetJsObject(),
		mpuCounter.GetJsObject(),
		mpwdCounter.GetJsObject(),
	}
	c.mu.Unlock()

	log.Printf("File parse and calculations took %s", time.Since(startTime))

	http.Redirect(w, r, "/stats?id="+id, 302)
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
	savedStats, ok := c.apiResponseCache[id]
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
