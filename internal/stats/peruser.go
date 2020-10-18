package stats

import (
	"github.com/psidex/MessengerStats/internal/messenger"
)

// MessagesPerUserJsObject is for putting the data in the structure that Highcharts requires.
type MessagesPerUserJsObject struct {
	Y    int    `json:"y"`
	Name string `json:"name"`
}

// MessagesPerUserCounter is for counting how many messages are sent per user.
// Pointers don't need to be used anywhere as the only field is a map which itself is a reference type.
type MessagesPerUserCounter struct {
	MessagesPerUser map[string]int
}

// NewMessagesPerUserCounter creates a new MessagesPerUserCounter.
func NewMessagesPerUserCounter() MessagesPerUserCounter {
	m := MessagesPerUserCounter{}
	m.MessagesPerUser = make(map[string]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerUserCounter) Update(message messenger.Message) {
	currentUser := message.SenderName
	if _, ok := m.MessagesPerUser[currentUser]; ok {
		m.MessagesPerUser[currentUser]++
	} else {
		m.MessagesPerUser[currentUser] = 1
	}
}

// GetJsObject returns the MessagesPerUserJsObject for passing to Highcharts.
func (m MessagesPerUserCounter) GetJsObject() []MessagesPerUserJsObject {
	var objects []MessagesPerUserJsObject
	for user, count := range m.MessagesPerUser {
		objects = append(objects, MessagesPerUserJsObject{
			Name: user,
			Y:    count,
		})
	}
	return objects
}
