package instagram

type Sender struct {
	Id string `json:"id"`
}

type Recipient struct {
	Id string `json:"id"`
}

type Payload struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Sticker_id    int64  `json:"sticker_id"`
	Reel_video_id string `json:"reel_video_id"`
}

type Attachment struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

type QuickReply struct {
	Payload string `json:"paylod"`
}

type ReplyTo struct {
	Mid string `json:"mid"`
}

type Message struct {
	Mid         string       `json:"mid"`
	Text        string       `json:"text"`
	Quick_reply QuickReply   `json:"quick_reply"`
	Reply_to    ReplyTo      `json:"reply_to"`
	Attachments []Attachment `json:"attachments"`
}

type MessageMetaData struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

type MessageWebhookObject struct {
	Object string         `json:"object"`
	Entry  []MessageEntry `json:"entry"`
}

type MessageEntry struct {
	Id        string            `json:"id"`
	Time      int64             `json:"time"`
	Messaging []MessageMetaData `json:"messaging"`
}
