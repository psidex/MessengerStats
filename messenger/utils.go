// Utilities for dealing with Messenger stuff.

package messenger

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// LoadMessengerJson takes a reader interface and a Messages struct, and unmarshals the data from the if into the struct.
func LoadMessengerJson(jsonFile io.Reader, messages *Messages) error {
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err := json.Unmarshal(byteValue, messages)
	return err
}
