package mojo

// Message is the base type for HTTP messages (Request and Response)
type Message struct {
	Headers Headers
	Content Asset
}

// NewMessage returns an initialized Message with empty content
func NewMessage() *Message {
	return &Message{Headers: Headers{}, Content: NewAsset("")}
}
