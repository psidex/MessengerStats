// ALl of the structs to unmarshal a JSON file of Messenger messages

package main

type nameObj struct {
	Name string `json:"name"`
}

type uriObj struct {
	Uri string `json:"uri"`
}

type uriWithTimestampObj struct {
	uriObj
	CreationTimestamp int64 `json:"creation_timestamp"`
}

type video struct {
	uriWithTimestampObj
	Thumbnail uriObj `json:"thumbnail"`
}

type share struct {
	Link string `json:"link"`
}

type reaction struct {
	Reaction string `json:"reaction"`
	Actor    string `json:"actor"`
}

type plan struct {
	Title     string `json:"title"`
	Location  string `json:"location"`
	Timestamp int64  `json:"int64"`
}

type Message struct {
	SenderName  string                `json:"sender_name"`
	TimestampMs int64                 `json:"timestamp_ms"`
	Sticker     uriObj                `json:"sticker"`
	Files       []uriWithTimestampObj `json:"files"`
	Photos      []uriWithTimestampObj `json:"photos"`
	Videos      []video               `json:"videos"`
	AudioFiles  []uriWithTimestampObj `json:"audio_files"`
	Gifs        []uriObj              `json:"gifs"`
	Content     string                `json:"content"`
	Share       share                 `json:"share"`
	Reactions   []reaction            `json:"reactions"`
	Plan        plan                  `json:"plan"`
	Type        string                `json:"type"`
	Users       []nameObj             `json:"users"`
}

type Messages struct {
	Participants       []nameObj `json:"participants"`
	Messages           []Message `json:"messages"`
	Title              string    `json:"title"`
	IsStillParticipant bool      `json:"is_still_participant"`
	ThreadType         string    `json:"thread_type"`
	ThreadPath         string    `json:"thread_path"`
}
