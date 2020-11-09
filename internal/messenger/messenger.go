package messenger

// These structs (courtesy of https://mholt.github.io/json-to-go/) are to be unmarshaled into using easyjson.
// If a field is required that is commented, uncomment it and re-run easyjson: "easyjson -all ./messenger.go".

// Message defines the structure of a single Messenger message.
type Message struct {
	SenderName  string `json:"sender_name"`
	TimestampMs int64  `json:"timestamp_ms"`
	//Content     string `json:"content,omitempty"`
	//Type        string `json:"type"`
	//Photos      []struct {
	//	URI               string `json:"uri"`
	//	CreationTimestamp int    `json:"creation_timestamp"`
	//} `json:"photos,omitempty"`
	//Share struct {
	//	Link      string `json:"link"`
	//	ShareText string `json:"share_text"`
	//} `json:"share,omitempty"`
	//Gifs []struct {
	//	URI string `json:"uri"`
	//} `json:"gifs,omitempty"`
	//Videos []struct {
	//	URI               string `json:"uri"`
	//	CreationTimestamp int    `json:"creation_timestamp"`
	//	Thumbnail         struct {
	//		URI string `json:"uri"`
	//	} `json:"thumbnail"`
	//} `json:"videos,omitempty"`
	//Reactions []struct {
	//	Reaction string `json:"reaction"`
	//	Actor    string `json:"actor"`
	//} `json:"reactions,omitempty"`
	//IP      string `json:"ip,omitempty"`
	//Sticker struct {
	//	URI string `json:"uri"`
	//} `json:"sticker,omitempty"`
	//Files []struct {
	//	URI               string `json:"uri"`
	//	CreationTimestamp int    `json:"creation_timestamp"`
	//} `json:"files,omitempty"`
	//CallDuration int  `json:"call_duration,omitempty"`
	//Missed       bool `json:"missed,omitempty"`
}

// Conversation defines the structure of a conversation in Messenger.
type Conversation struct {
	//Participants []struct {
	//	Name string `json:"name"`
	//} `json:"participants"`
	Messages []Message `json:"messages"`
	Title    string    `json:"title"`
	//IsStillParticipant bool      `json:"is_still_participant"`
	//ThreadType         string    `json:"thread_type"`
	//ThreadPath         string    `json:"thread_path"`
}

// NewConversation attempts to unmarshal the data from the byte slice into a new Conversation struct.
func NewConversation(data []byte) (*Conversation, error) {
	conv := &Conversation{}
	if err := conv.UnmarshalJSON(data); err != nil {
		return conv, err
	}
	return conv, nil
}
