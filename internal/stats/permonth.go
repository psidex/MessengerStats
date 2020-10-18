package stats

import (
	"fmt"
	"github.com/psidex/MessengerStats/internal/messenger"
	"sort"
	"time"
)

// MessagesPerMonthJsObject is for putting the data in the structure that Highcharts requires.
type MessagesPerMonthJsObject struct {
	Data       []int    `json:"data"`
	Categories []string `json:"categories"`
}

// MessagesPerMonthCounter is for counting how many messages are sent per month.
type MessagesPerMonthCounter struct {
	messagesPerYearMonth map[int]map[int]int
}

// NewMessagesPerMonthCounter creates a new MessagesPerMonthCounter.
func NewMessagesPerMonthCounter() MessagesPerMonthCounter {
	m := MessagesPerMonthCounter{}
	m.messagesPerYearMonth = make(map[int]map[int]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerMonthCounter) Update(message messenger.Message) {
	timestamp := message.TimestampMs / 1000
	year, mo, _ := time.Unix(timestamp, 0).Date()
	month := int(mo)

	if _, ok := m.messagesPerYearMonth[year]; ok {
		if _, ok := m.messagesPerYearMonth[year][month]; ok {
			m.messagesPerYearMonth[year][month]++
		} else {
			m.messagesPerYearMonth[year][month] = 1
		}
	} else {
		m.messagesPerYearMonth[year] = map[int]int{month: 1}
	}
}

// GetJsObject returns the MessagesPerMonthJsObject for passing to Highcharts.
func (m MessagesPerMonthCounter) GetJsObject() MessagesPerMonthJsObject {
	messagesPerYearMonthSortedKeys := make([]int, len(m.messagesPerYearMonth))
	i := 0
	for k := range m.messagesPerYearMonth {
		messagesPerYearMonthSortedKeys[i] = k
		i++
	}
	sort.Ints(messagesPerYearMonthSortedKeys)

	firstYear := messagesPerYearMonthSortedKeys[0]
	currentYear := time.Now().Year()

	// Fill out our data so that months with 0 messages actually show 0.
	// TODO: Will set all months to 0 even if not happened yet (important for first and last year)
	// From the first year to current year.
	for iterYear := firstYear; iterYear <= currentYear; iterYear++ {
		// If there is data for that year.
		if _, ok := m.messagesPerYearMonth[iterYear]; ok {
			// For each month.
			for iterMonth := 1; iterMonth <= 12; iterMonth++ {
				// If there is not data, set the value to 0 (before there would have been no such key).
				if _, ok := m.messagesPerYearMonth[iterYear][iterMonth]; !ok {
					m.messagesPerYearMonth[iterYear][iterMonth] = 0
				}
			}
		} else {
			// Possible if conversation breaks for a year or more.
			m.messagesPerYearMonth[iterYear] = map[int]int{
				1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0,
			}
		}
	}

	// Construct object.
	obj := MessagesPerMonthJsObject{}
	for _, year := range messagesPerYearMonthSortedKeys {
		monthMap := m.messagesPerYearMonth[year]
		for _, month := range sortedMapKeys(monthMap) {
			text := fmt.Sprintf("%d-%d", year, month)
			obj.Categories = append(obj.Categories, text)
			obj.Data = append(obj.Data, m.messagesPerYearMonth[year][month])
		}
	}
	return obj
}

// sortedMapKeys takes a map and returns the sorted keys.
func sortedMapKeys(mapping map[int]int) []int {
	var keys []int
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
