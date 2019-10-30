// For dealing with gathering statistics from a Message struct.

package main

import (
	"github.com/bbalet/stopwords"
	"strings"
	"time"
	"unicode"
)

type date struct {
	Year  int
	Month time.Month
}

// getSentDate takes a Message struct and returns a date struct of the date that it was sent.
func getSentDate(message Message) date {
	// FB uses ms instead of s
	ts := message.Timestamp / 1000
	y, m, _ := time.Unix(ts, 0).Date()
	return date{y, m}
}

// CountMessagesPerDate takes a Messages struct and counts how many messages were sent on each date.
// The date is represented as a Year + Month, not including the day.
func CountMessagesPerDate(messages *Messages) map[date]int {
	messagesPerDate := make(map[date]int)

	for _, message := range messages.Messages {

		currentDate := getSentDate(message)
		// https://stackoverflow.com/a/2050629/6396652
		if _, ok := messagesPerDate[currentDate]; ok {
			messagesPerDate[currentDate]++
		} else {
			messagesPerDate[currentDate] = 1
		}

	}

	return messagesPerDate
}

// CountMessagesPerUser takes a Messages struct and counts how many messages were sent by each user.
func CountMessagesPerUser(messages *Messages) map[string]int {
	messagesPerUser := make(map[string]int)

	for _, message := range messages.Messages {

		currentUser := message.SenderName
		if _, ok := messagesPerUser[currentUser]; ok {
			messagesPerUser[currentUser]++
		} else {
			messagesPerUser[currentUser] = 1
		}

	}

	return messagesPerUser
}

// Returns true only if the string contains all ASCII characters.
// From https://stackoverflow.com/a/53069799/6396652.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// CountWords takes a Messages struct and counts how many times each word is said. Only counts ASCII words.
func CountWords(messages *Messages) map[string]int {
	wordCount := make(map[string]int)

	for _, message := range messages.Messages {

		formattedContent := strings.ReplaceAll(message.Content, "\n", " ")
		words := strings.Split(formattedContent, " ")

		for _, word := range words {

			if !isASCII(word) {
				continue
			}

			// Format and clean string.
			formattedWord := word
			formattedWord = stopwords.CleanString(formattedWord, "en", false)
			formattedWord = strings.ToLower(formattedWord)
			formattedWord = strings.TrimSpace(formattedWord)

			if _, ok := wordCount[formattedWord]; ok {
				wordCount[formattedWord]++
			} else {
				wordCount[formattedWord] = 1
			}

		}

	}

	return wordCount
}
