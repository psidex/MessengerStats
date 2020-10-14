package stats

import (
	"github.com/psidex/MessengerStats/internal/messenger"
	"time"
)

// MessagesPerWeekdayCounter is for counting how many messages are sent per weekday.
type MessagesPerWeekdayCounter struct {
	MessagesPerDay map[string]int
}

// NewMessagesPerWeekdayCounter creates a new MessagesPerWeekdayCounter.
func NewMessagesPerWeekdayCounter() MessagesPerWeekdayCounter {
	m := MessagesPerWeekdayCounter{}
	m.MessagesPerDay = make(map[string]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerWeekdayCounter) Update(message messenger.Message) {
	ts := message.TimestampMs / 1000
	currentDay := time.Unix(ts, 0).Weekday().String()

	if _, ok := m.MessagesPerDay[currentDay]; ok {
		m.MessagesPerDay[currentDay]++
	} else {
		m.MessagesPerDay[currentDay] = 1
	}
}
