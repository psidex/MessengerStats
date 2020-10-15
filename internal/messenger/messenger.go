package messenger

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Message defines the structure of a single Messenger message.
type Message struct {
	SenderName  string `json:"sender_name"`
	TimestampMs int64  `json:"timestamp_ms"`
	Content     string `json:"content,omitempty"`
	Type        string `json:"type"`
	Photos      []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"photos,omitempty"`
	Share struct {
		Link      string `json:"link"`
		ShareText string `json:"share_text"`
	} `json:"share,omitempty"`
	Gifs []struct {
		URI string `json:"uri"`
	} `json:"gifs,omitempty"`
	Videos []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
		Thumbnail         struct {
			URI string `json:"uri"`
		} `json:"thumbnail"`
	} `json:"videos,omitempty"`
	Reactions []struct {
		Reaction string `json:"reaction"`
		Actor    string `json:"actor"`
	} `json:"reactions,omitempty"`
	IP      string `json:"ip,omitempty"`
	Sticker struct {
		URI string `json:"uri"`
	} `json:"sticker,omitempty"`
	Files []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"files,omitempty"`
	CallDuration int  `json:"call_duration,omitempty"`
	Missed       bool `json:"missed,omitempty"`
}

// Conversation defines the structure of a conversation in Messenger.
type Conversation struct {
	Participants []struct {
		Name string `json:"name"`
	} `json:"participants"`
	Messages           []Message `json:"messages"`
	Title              string    `json:"title"`
	IsStillParticipant bool      `json:"is_still_participant"`
	ThreadType         string    `json:"thread_type"`
	ThreadPath         string    `json:"thread_path"`
}

// NewConversation attempts to unmarshal the data from the io.Reader into a new Conversation struct.
func NewConversation(reader io.Reader) (*Conversation, error) {
	conv := &Conversation{}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return &Conversation{}, err
	}
	err = json.Unmarshal(bytes, conv)
	if err != nil {
		return &Conversation{}, err
	}
	return conv, nil
}
