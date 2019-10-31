// For dealing with gathering statistics from a Message struct.

package messenger

import (
	"github.com/bbalet/stopwords"
	"strings"
	"time"
	"unicode"
)

// getYearMonth returns the year and month the message was sent as integers (e.g. 2000, 12 for Dec 2000).
func getYearMonth(message Message) (int, int) {
	// FB uses ms instead of s
	ts := message.Timestamp / 1000
	y, m, _ := time.Unix(ts, 0).Date()
	return y, int(m)
}

// CountMessagesPerMonth takes a Messages struct and counts how many messages were sent on each month, categorized by year.
// Example map using data from Nov and Dec in 2017: {2017: {11: 69, 12: 420}}
func CountMessagesPerMonth(messages *Messages) map[int]map[int]int {
	messagesPerDate := make(map[int]map[int]int)

	for _, message := range messages.Messages {

		year, month := getYearMonth(message)

		// https://stackoverflow.com/a/2050629/6396652
		if _, ok := messagesPerDate[year]; ok {
			if _, ok := messagesPerDate[year][month]; ok {
				messagesPerDate[year][month]++
			} else {
				messagesPerDate[year][month] = 1
			}
		} else {
			messagesPerDate[year] = map[int]int{month: 1}
		}

	}

	return messagesPerDate
}

// getSentDay returns the weekday that the message was sent on.
func getSentDay(message Message) string {
	ts := message.Timestamp / 1000
	return time.Unix(ts, 0).Weekday().String()
}

// CountMessagesPerWeekday takes a Messages struct and counts how many messages were sent on each day.
func CountMessagesPerWeekday(messages *Messages) map[string]int {
	messagesPerDate := make(map[string]int)

	for _, message := range messages.Messages {

		currentDay := getSentDay(message)
		if _, ok := messagesPerDate[currentDay]; ok {
			messagesPerDate[currentDay]++
		} else {
			messagesPerDate[currentDay] = 1
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
