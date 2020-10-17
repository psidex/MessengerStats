package stats

import (
	"github.com/bbalet/stopwords"
	"github.com/psidex/MessengerStats/internal/messenger"
	"strings"
	"unicode"
)

// WordCountJsObject is for putting the data in the structure that Highcharts requires.
type WordCountJsObject struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

// WordCounter is for counting words.
type WordCounter struct {
	counts map[string]int
}

// NewWordCounter creates a new WordCounter.
func NewWordCounter() WordCounter {
	m := WordCounter{}
	m.counts = make(map[string]int)
	return m
}

// Update updates the counter with data from a single message.
func (w WordCounter) Update(message messenger.Message) {
	// Remove all stop words, split line into words, only process non-blank and ascii text.
	contentBytes := []byte(message.Content)
	wordBytes := stopwords.Clean(contentBytes, "en", false)
	for _, word := range strings.Split(string(wordBytes), " ") {
		if word == "" || !isASCII(word) {
			continue
		}
		if _, ok := w.counts[word]; ok {
			w.counts[word]++
		} else {
			w.counts[word] = 1
		}
	}
}

// GetJsObject returns the array of WordCountJsObject for passing to Highcharts.
func (w WordCounter) GetJsObject() []WordCountJsObject {
	var objects []WordCountJsObject
	for word, count := range w.counts {
		if count <= 100 {
			// TODO: Maybe this should be variable.
			continue
		}
		objects = append(objects, WordCountJsObject{
			Name:   word,
			Weight: count,
		})
	}
	return objects
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
