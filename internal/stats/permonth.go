package stats

import (
	"github.com/psidex/MessengerStats/internal/messenger"
	"time"
)

// TODO: If a year and/or month don't have any entries, set to 0 so it still shows on chart.

// MessagesPerMonthCounter is for counting how many messages are sent per month.
type MessagesPerMonthCounter struct {
	MessagesPerYearMonth map[int]map[int]int
}

// NewMessagesPerMonthCounter creates a new MessagesPerMonthCounter.
func NewMessagesPerMonthCounter() MessagesPerMonthCounter {
	m := MessagesPerMonthCounter{}
	m.MessagesPerYearMonth = make(map[int]map[int]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerMonthCounter) Update(message messenger.Message) {
	timestamp := message.TimestampMs / 1000
	year, mo, _ := time.Unix(timestamp, 0).Date()
	month := int(mo)

	if _, ok := m.MessagesPerYearMonth[year]; ok {
		if _, ok := m.MessagesPerYearMonth[year][month]; ok {
			m.MessagesPerYearMonth[year][month]++
		} else {
			m.MessagesPerYearMonth[year][month] = 1
		}
	} else {
		m.MessagesPerYearMonth[year] = map[int]int{month: 1}
	}
}
