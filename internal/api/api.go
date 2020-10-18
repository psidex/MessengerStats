package api

import (
	"encoding/json"
	"github.com/psidex/MessengerStats/internal/messenger"
	"github.com/psidex/MessengerStats/internal/stats"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"sync"
	"time"
)

// statsApiResponse holds all the data to be returned by the /stats endpoint.
type statsApiResponse struct {
	singleUse          bool                             // If true, erase from cache after request.
	Error              string                           `json:"error"`
	Title              string                           `json:"conversation_title"`
	MessagesPerMonth   stats.MessagesPerMonthJsObject   `json:"messages_per_month"`
	MessagesPerUser    []stats.MessagesPerUserJsObject  `json:"messages_per_user"`
	MessagesPerWeekday stats.MessagesPerWeekdayJsObject `json:"messages_per_weekday"`
}

// uploadApiResponse holds all the data to be returned by the /upload endpoint.
type uploadApiResponse struct {
	Error string `json:"error"`
	Id    string `json:"id"`
}

// StatsApi contains the data and functions for generating conversation statistics.
type StatsApi struct {
	apiResponseCache map[string]statsApiResponse
	mu               *sync.Mutex // Controls access to apiResponseCache
}

func NewStatsApi() StatsApi {
	c := StatsApi{}
	c.mu = &sync.Mutex{}
	c.apiResponseCache = make(map[string]statsApiResponse)
	return c
}

// UploadHandler is a HTTP handler for the /api/upload endpoint.
func (s StatsApi) UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseEncoder := json.NewEncoder(w)

	const _500MB = 1024 * 1024 * 500 // Use up to 500MB of RAM, rest goes on disk.
	if err := r.ParseMultipartForm(_500MB); nil != err {
		_ = responseEncoder.Encode(uploadApiResponse{
			Error: err.Error(),
		})
		return
	}

	headers, ok := r.MultipartForm.File["messenger_files"]
	if !ok {
		_ = responseEncoder.Encode(uploadApiResponse{
			Error: "could not find messenger_files in form",
		})
		return
	}

	mpuCounter := stats.NewMessagesPerUserCounter()
	mpwdCounter := stats.NewMessagesPerWeekdayCounter()
	mpmCounter := stats.NewMessagesPerMonthCounter()

	title := ""
	id := ksuid.New().String()
	startTime := time.Now()

	for _, header := range headers {
		file, err := header.Open()

		if err != nil {
			_ = responseEncoder.Encode(uploadApiResponse{
				Error: err.Error(),
			})
			return
		}

		conversation, err := messenger.NewConversation(file)
		_ = file.Close()
		if err != nil {
			_ = responseEncoder.Encode(uploadApiResponse{
				Error: err.Error(),
			})
			return
		}

		title = conversation.Title

		for _, message := range conversation.Messages {
			mpuCounter.Update(message)
			mpwdCounter.Update(message)
			mpmCounter.Update(message)
		}
	}

	cacheParam := r.URL.Query().Get("cache")
	singleUse := true
	if cacheParam == "true" {
		singleUse = false
	}

	s.mu.Lock()
	s.apiResponseCache[id] = statsApiResponse{
		singleUse:          singleUse,
		Error:              "",
		Title:              title,
		MessagesPerMonth:   mpmCounter.GetJsObject(),
		MessagesPerUser:    mpuCounter.GetJsObject(),
		MessagesPerWeekday: mpwdCounter.GetJsObject(),
	}
	s.mu.Unlock()

	log.Printf("File parse and calculations took %s", time.Since(startTime))

	_ = responseEncoder.Encode(uploadApiResponse{
		Error: "",
		Id:    id,
	})
}

// StatsHandler is a HTTP handler that takes an ID (param in a url) and returns the respective data from `c`.
func (s StatsApi) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	responseEncoder := json.NewEncoder(w)

	id := r.URL.Query().Get("id")
	if id == "" {
		_ = responseEncoder.Encode(statsApiResponse{
			Error: "no \"id\" parameter present in request",
		})
		return
	}

	log.Println("Stats for ID requested:", id)

	s.mu.Lock()
	savedStats, ok := s.apiResponseCache[id]
	if ok && savedStats.singleUse {
		delete(s.apiResponseCache, id)
	}
	s.mu.Unlock()

	if !ok {
		_ = responseEncoder.Encode(statsApiResponse{
			Error: "invalid ID",
		})
		return
	}

	_ = responseEncoder.Encode(savedStats)
}
