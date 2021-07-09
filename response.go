package mojo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response represents a response to the client.
type Response struct {
	Message
	Writer http.ResponseWriter
	Code   int
	Status string
	raw    *http.Response
}

// NewResponse returns a new, empty response with sensible defaults.
func NewResponse(opts ...interface{}) *Response {
	res := &Response{Message: Message{Headers: Headers{}, Content: NewAsset("")}}
	for _, opt := range opts {
		switch v := opt.(type) {
		case http.ResponseWriter:
			res.Writer = v
		default:
			panic(fmt.Sprintf("Unknown option type %T\n", v))
		}
	}
	return res
}

// Read populates this response from the given http.Response
func (res *Response) Read(raw *http.Response) {
	res.raw = raw
	res.Headers = Headers(raw.Header)
	res.Content = NewAsset(raw.Body)
}

// JSON encodes the given argument as JSON and updates the response's
// Content-Type header.
func (res *Response) JSON(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	res.Content = NewAsset(json)
	res.Headers["Content-Type"] = []string{"application/json"}
}

// Text sets the response's content and updates the Content-Type header
// to "text/plain"
func (res *Response) Text(str string) {
	res.Content = NewAsset(str)
	res.Headers["Content-Type"] = []string{"text/plain"}
}
