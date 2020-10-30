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
	messagesPerMonth map[int]map[int]int
}

// NewMessagesPerMonthCounter creates a new MessagesPerMonthCounter.
func NewMessagesPerMonthCounter() MessagesPerMonthCounter {
	m := MessagesPerMonthCounter{}
	m.messagesPerMonth = make(map[int]map[int]int)
	return m
}

// Update updates the counter with data from a single message.
func (m MessagesPerMonthCounter) Update(message messenger.Message) {
	timestamp := message.TimestampMs / 1000
	year, mo, _ := time.Unix(timestamp, 0).Date()
	month := int(mo)

	if _, ok := m.messagesPerMonth[year]; ok {
		if _, ok := m.messagesPerMonth[year][month]; ok {
			m.messagesPerMonth[year][month]++
		} else {
			m.messagesPerMonth[year][month] = 1
		}
	} else {
		m.messagesPerMonth[year] = map[int]int{month: 1}
	}
}

// GetJsObject returns the MessagesPerMonthJsObject for passing to Highcharts.
func (m MessagesPerMonthCounter) GetJsObject() MessagesPerMonthJsObject {
	if len(m.messagesPerMonth) <= 0 {
		return MessagesPerMonthJsObject{}
	}

	var messagesPerMonthSortedKeys []int
	for k := range m.messagesPerMonth {
		messagesPerMonthSortedKeys = append(messagesPerMonthSortedKeys, k)
	}
	sort.Ints(messagesPerMonthSortedKeys)

	firstYear := messagesPerMonthSortedKeys[0]
	// A year is never created without setting at least 1 month so we don't need to worry about a map with no values.
	firstMonth := sortedMapKeys(m.messagesPerMonth[firstYear])[0]
	currentYear := time.Now().Year()
	currentMonth := int(time.Now().Month())

	// Fill out our data so that months with 0 messages actually show 0.
	for rd := rangeMonths(firstYear, firstMonth, currentYear, currentMonth); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		year, mo, _ := date.Date()
		month := int(mo)
		if _, ok := m.messagesPerMonth[year][month]; !ok {
			m.messagesPerMonth[year][month] = 0
		}
	}

	// Construct object.
	obj := MessagesPerMonthJsObject{}
	for _, year := range messagesPerMonthSortedKeys {
		monthMap := m.messagesPerMonth[year]
		for _, month := range sortedMapKeys(monthMap) {
			text := fmt.Sprintf("%d-%d", year, month)
			obj.Categories = append(obj.Categories, text)
			obj.Data = append(obj.Data, m.messagesPerMonth[year][month])
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

// rangeMonths returns a range function over start date to end date inclusive, returning incrementing months.
// After the end of the range, the range function returns a zero date (date.IsZero() is true).
func rangeMonths(startYear, startMonth, endYear, endMonth int) func() time.Time {
	// Edited from https://stackoverflow.com/a/50989054/6396652
	start := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 1, 0)
		return date
	}
}
