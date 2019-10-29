package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// loadJson takes a path to a json file and a struct, and unmarshals the data from the file into the struct
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

	messagesPerDate := CountMessagesPerDate(messages)
	for d, messageCount := range messagesPerDate {
		log.Println("Date:", d, "Message count:", messageCount)
	}

	log.Println("Done. Messages processed:", len(messages.Messages))
}
