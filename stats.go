// For dealing with gathering statistics from a Message struct

package main

import "time"

type date struct {
	Year  int
	Month time.Month
	Day   int
}

// getSentDate takes a Message struct and returns a date struct of the date that it was sent
func getSentDate(message Message) date {
	// FB uses ms instead of s
	ts := message.TimestampMs / 1000
	y, m, d := time.Unix(ts, 0).Date()
	return date{y, m, d}
}

// CountMessagesPerDate tales a Message struct and counts how many messages were sent on each date
func CountMessagesPerDate(messages *Messages) map[date]int {
	messagesPerDate := make(map[date]int)

	for i := 0; i < len(messages.Messages); i++ {

		currentDate := getSentDate(messages.Messages[i])
		// https://stackoverflow.com/a/2050629/6396652
		if _, ok := messagesPerDate[currentDate]; ok {
			messagesPerDate[currentDate]++
		} else {
			messagesPerDate[currentDate] = 1
		}

	}

	return messagesPerDate
}
