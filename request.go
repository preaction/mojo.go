package mojo

import (
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
