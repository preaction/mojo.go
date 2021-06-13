package mojo

// Message is the base type for HTTP messages (Request and Response)
type Message struct {
	Headers Headers
	Content string
}
