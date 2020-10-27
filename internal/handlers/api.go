package handlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/psidex/MessengerStats/internal/messenger"
	"github.com/psidex/MessengerStats/internal/stats"
	"log"
	"net/http"
	"time"
)

var upgrader websocket.Upgrader

// apiResponse holds all the data to be returned by API endpoint.
type apiResponse struct {
	Error              string                           `json:"error"`
	Title              string                           `json:"conversation_title"`
	MessagesPerMonth   stats.MessagesPerMonthJsObject   `json:"messages_per_month"`
	MessagesPerUser    []stats.MessagesPerUserJsObject  `json:"messages_per_user"`
	MessagesPerWeekday stats.MessagesPerWeekdayJsObject `json:"messages_per_weekday"`
}

// WebSocketApi is a handler for the api's websocket endpoint.
func WebSocketApi(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrader error: %s\n", err)
		return
	}
	defer ws.Close()

	t, fileCountBytes, _ := ws.ReadMessage()
	if t != websocket.BinaryMessage || len(fileCountBytes) != 1 {
		js, _ := json.Marshal(apiResponse{Error: "file count should be a single byte"})
		_ = ws.WriteMessage(websocket.BinaryMessage, js)
		return
	}
	fileCount := int(fileCountBytes[0])

	if fileCount <= 0 {
		js, _ := json.Marshal(apiResponse{Error: "file count must be > 0"})
		_ = ws.WriteMessage(websocket.BinaryMessage, js)
		return
	}

	mpuCounter := stats.NewMessagesPerUserCounter()
	mpwdCounter := stats.NewMessagesPerWeekdayCounter()
	mpmCounter := stats.NewMessagesPerMonthCounter()

	title := ""
	startTime := time.Now()

	// Receive each file.
	for i := 0; i < fileCount; i++ {
		fileReadTimeStart := time.Now()
		t, fileBytes, _ := ws.ReadMessage()
		log.Printf("File read %d took %s\n", i, time.Since(fileReadTimeStart))

		if t != websocket.BinaryMessage {
			js, _ := json.Marshal(apiResponse{Error: "files should be sent as binary data"})
			_ = ws.WriteMessage(websocket.BinaryMessage, js)
			return
		}

		conversation, err := messenger.NewConversation(fileBytes)
		if err != nil {
			js, _ := json.Marshal(apiResponse{Error: err.Error()})
			_ = ws.WriteMessage(websocket.BinaryMessage, js)
			return
		}

		title = conversation.Title

		for _, message := range conversation.Messages {
			mpuCounter.Update(message)
			mpwdCounter.Update(message)
			mpmCounter.Update(message)
		}
	}

	log.Printf("File processing took %s\n", time.Since(startTime))

	jsonResponse, err := json.Marshal(apiResponse{
		Title:              title,
		MessagesPerMonth:   mpmCounter.GetJsObject(),
		MessagesPerUser:    mpuCounter.GetJsObject(),
		MessagesPerWeekday: mpwdCounter.GetJsObject(),
	})
	if err != nil {
		log.Printf("Error marshalling apiResponse: %s\n", err)
		return
	}

	err = ws.WriteMessage(websocket.BinaryMessage, jsonResponse)
	if err != nil {
		log.Printf("Error sending apiResponse: %s\n", err)
		return
	}
}
