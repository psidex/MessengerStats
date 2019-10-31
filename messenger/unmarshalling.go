// ALl of the structs to unmarshal a JSON file of Messenger messages.

package messenger

type NameObj struct {
	Name string `json:"name"`
}

type UriObj struct {
	Uri string `json:"uri"`
}

type UriWithTimestampObj struct {
	UriObj
	CreationTimestamp int64 `json:"creation_timestamp"`
}

type Video struct {
	UriWithTimestampObj
	Thumbnail UriObj `json:"thumbnail"`
}

type Share struct {
	Link string `json:"link"`
}

type Reaction struct {
	Reaction string `json:"reaction"`
	Actor    string `json:"actor"`
}

type Plan struct {
	Title     string `json:"title"`
	Location  string `json:"location"`
	Timestamp int64  `json:"timestamp"`
}

type Message struct {
	SenderName string                `json:"sender_name"`
	Timestamp  int64                 `json:"timestamp_ms"`
	Sticker    UriObj                `json:"sticker"`
	Files      []UriWithTimestampObj `json:"files"`
	Photos     []UriWithTimestampObj `json:"photos"`
	Videos     []Video               `json:"videos"`
	AudioFiles []UriWithTimestampObj `json:"audio_files"`
	Gifs       []UriObj              `json:"gifs"`
	Content    string                `json:"content"`
	Share      Share                 `json:"Share"`
	Reactions  []Reaction            `json:"reactions"`
	Plan       Plan                  `json:"Plan"`
	Type       string                `json:"type"`
	Users      []NameObj             `json:"users"`
}

type Messages struct {
	Participants       []NameObj `json:"participants"`
	Messages           []Message `json:"messages"`
	Title              string    `json:"title"`
	IsStillParticipant bool      `json:"is_still_participant"`
	ThreadType         string    `json:"thread_type"`
	ThreadPath         string    `json:"thread_path"`
}
