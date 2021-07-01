package mojo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request represents an HTTP request
type Request struct {
	Message
	Method      string
	URL         *url.URL
	Params      Parameters
	QueryParams Parameters
	BodyParams  Parameters

	raw *http.Request
}

// NewRequest builds a new request object
func NewRequest(method string, inputURL string) *Request {
	// XXX: Add body optional argument and auto-populate query or body params as needed
	requestURL, err := url.ParseRequestURI(inputURL)
	if err != nil {
		panic(fmt.Sprintf("Could not parse URL: %v", err))
	}
	return &Request{Method: method, URL: requestURL}

}

// Read populates this request from the given http.Request
func (req *Request) Read(raw *http.Request) {
	req.raw = raw
	req.URL = raw.URL
	req.Method = raw.Method

	req.Params = Parameters{}
	req.QueryParams = Parameters{}
	for k, v := range raw.URL.Query() {
		req.QueryParams[k] = v
		req.Params[k] = v
	}

	req.BodyParams = Parameters{}
	raw.ParseForm()
	for k, v := range raw.PostForm {
		req.BodyParams[k] = v
		req.Params[k] = v
	}

	req.Headers = Headers(raw.Header)
	// The Host header was removed by Go, so we have to put it back
	req.Headers["Host"] = []string{raw.Host}
}

// Param gets the first value for the given parameter. Body parameters (POST
// forms) take precedence over query parameters (URLs). To get every
// value, see EveryParam.
func (req *Request) Param(name string) string {
	return req.Params.Param(name)
}

// EveryParam gets all values for the given parameter. Body parameters (POST
// forms) take precedence over query parameters (URLs).
func (req *Request) EveryParam(name string) []string {
	return req.Params.EveryParam(name)
}

// readBody reads the request body if necessary, caches it in the
// Request object, and returns it.
func (req *Request) readContent() string {
	if req.Content != "" {
		return req.Content
	}
	body, err := io.ReadAll(req.raw.Body)
	if err != nil {
		panic(err)
	}
	req.Content = string(body)
	return req.Content
}

// JSON reads the request body and unmarshals into the given type
// pointer. Returns an error if JSON parsing fails.
func (req *Request) JSON(empty interface{}) error {
	content := req.readContent()
	err := json.Unmarshal([]byte(content), empty)
	return err
}
