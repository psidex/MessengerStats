package stats

import (
	"github.com/psidex/MessengerStats/internal/messenger"
	"time"
)

// MessagesPerWeekdayJsObject is for putting the data in the structure that Highcharts requires.
type MessagesPerWeekdayJsObject struct {
	Data       []int    `json:"data"`
	Categories []string `json:"categories"`
}

// MessagesPerWeekdayCounter is for counting how many messages are sent per weekday.
type MessagesPerWeekdayCounter struct {
	messagesPerDay map[string]int
}

// NewMessagesPerWeekdayCounter creates a new MessagesPerWeekdayCounter.
func NewMessagesPerWeekdayCounter() MessagesPerWeekdayCounter {
	m := MessagesPerWeekdayCounter{}
	m.messagesPerDay = make(map[string]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerWeekdayCounter) Update(message messenger.Message) {
	ts := message.TimestampMs / 1000
	currentDay := time.Unix(ts, 0).Weekday().String()

	if _, ok := m.messagesPerDay[currentDay]; ok {
		m.messagesPerDay[currentDay]++
	} else {
		m.messagesPerDay[currentDay] = 1
	}
}

// GetJsObject returns the MessagesPerWeekdayJsObject for passing to Highcharts.
func (m MessagesPerWeekdayCounter) GetJsObject() MessagesPerWeekdayJsObject {
	obj := MessagesPerWeekdayJsObject{}
	// Make sure we create ordered arrays.
	orderedDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, weekday := range orderedDays {
		obj.Categories = append(obj.Categories, weekday)

		if messageCount, ok := m.messagesPerDay[weekday]; ok {
			obj.Data = append(obj.Data, messageCount)
		} else {
			obj.Data = append(obj.Data, 0)
		}

	}
	return obj
}
