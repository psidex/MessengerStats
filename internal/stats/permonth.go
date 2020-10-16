package stats

import (
	"github.com/psidex/MessengerStats/internal/messenger"
	"sort"
	"time"
)

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

// Finalize takes the calculated data and adds in months that should have a value of 0 but currently don't exist.
// This can't be done in Update because it's not guaranteed that it will be called on data in chronological order.
func (m MessagesPerMonthCounter) Finalize() {
	// Find the lowest year in the map.
	// https://stackoverflow.com/a/23332089/6396652
	keys := make([]int, len(m.MessagesPerYearMonth))
	i := 0
	for k := range m.MessagesPerYearMonth {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	firstYear := keys[0]
	currentYear := time.Now().Year()

	// TODO: Will set all months to 0 even if not happened yet (important for first and last year)
	// From the first year to current year.
	for iterYear := firstYear; iterYear <= currentYear; iterYear++ {
		// If there is data for that year.
		if _, ok := m.MessagesPerYearMonth[iterYear]; ok {
			// For each month.
			for iterMonth := 1; iterMonth <= 12; iterMonth++ {
				// If there is not data, set the value to 0 (before there would have been no such key).
				if _, ok := m.MessagesPerYearMonth[iterYear][iterMonth]; !ok {
					m.MessagesPerYearMonth[iterYear][iterMonth] = 0
				}
			}
		} else {
			// Possible if conversation breaks for a year or more.
			m.MessagesPerYearMonth[iterYear] = map[int]int{
				1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0,
			}
		}
	}
}
