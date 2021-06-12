package mojo

import (
	"encoding/json"
	"net/http"
)

// Response represents a response to the client.
type Response struct {
	Writer  http.ResponseWriter
	Code    int
	Message string
	Body    string
	Headers Headers
}

// NewResponse returns a new, empty response with sensible defaults.
func NewResponse() *Response {
	return &Response{Headers: Headers{}}
}

// JSON encodes the given argument as JSON and updates the response's
// Content-Type header.
func (res *Response) JSON(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	res.Body = string(json)
	res.Headers["Content-Type"] = []string{"application/json"}
}
