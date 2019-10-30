package main

import (
	"encoding/json"
	"github.com/psykhi/wordclouds"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// loadJson takes a path to a json file and a struct, and unmarshals the data from the file into the struct.
func loadJson(pathToFile string, messages *Messages) {
	jsonFile, err := os.Open(pathToFile)
	if err != nil {
		log.Fatal("Couldn't open messages file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, messages)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON data: ", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No messenger directory specified")
	}

	messagesDir := os.Args[1]
	messagesJsonFile := path.Join(messagesDir, "message_1.json")
	log.Println("Opening", messagesJsonFile)

	messages := &Messages{}
	loadJson(messagesJsonFile, messages)

	//messagesPerDate := CountMessagesPerDate(messages)
	//for d, messageCount := range messagesPerDate {
	//	log.Println("Date:", d, "Message count:", messageCount)
	//}
	//
	//messagesPerUser := CountMessagesPerUser(messages)
	//for user, messageCount := range messagesPerUser {
	//	log.Println("User:", user, "Message count:", messageCount)
	//}

	wordCount := CountWords(messages)
	newWordCount := make(map[string]int)
	for word, wordCount := range wordCount {
		if wordCount > 100 && len(word) > 1 {
			newWordCount[word] = wordCount
		}
	}

	colors := []color.Color{
		color.RGBA{0x1b, 0x1b, 0x1b, 0xff},
		color.RGBA{0x48, 0x48, 0x4B, 0xff},
		color.RGBA{0x59, 0x3a, 0xee, 0xff},
		color.RGBA{0x65, 0xCD, 0xFA, 0xff},
		color.RGBA{0x70, 0xD6, 0xBF, 0xff},
	}

	w := wordclouds.NewWordcloud(
		newWordCount,
		wordclouds.FontFile("Roboto-Regular.ttf"),
		wordclouds.FontMaxSize(700),
		wordclouds.FontMinSize(10),
		wordclouds.Colors(colors),
		wordclouds.Height(4096),
		wordclouds.Width(4096),
	)

	img := w.Draw()

	outputFile, err := os.Create("output.png")
	if err != nil {
		log.Fatal("Failed to open output.png:", err)
	}
	_ = png.Encode(outputFile, img)
	_ = outputFile.Close()

	log.Println("Done. Messages processed:", len(messages.Messages))
}
