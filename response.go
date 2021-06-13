package mojo

import (
	"encoding/json"
	"net/http"
)

// Response represents a response to the client.
type Response struct {
	Message
	Writer http.ResponseWriter
	Code   int
	Status string
}

// NewResponse returns a new, empty response with sensible defaults.
func NewResponse() *Response {
	return &Response{Message: Message{Headers: Headers{}}}
}

// JSON encodes the given argument as JSON and updates the response's
// Content-Type header.
func (res *Response) JSON(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	res.Content = string(json)
	res.Headers["Content-Type"] = []string{"application/json"}
}
