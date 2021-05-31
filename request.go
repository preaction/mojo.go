package mojo

import (
	"net/http"
	"net/url"
)

type Request struct {
	Message
	Method string
	URL    *url.URL
	raw    *http.Request
}

func (req *Request) Read(raw *http.Request) {
	req.raw = raw
	req.URL = raw.URL
	req.Method = raw.Method
}

func (req *Request) Param(name string) string {
	// XXX: Create mojo.Parameters
	// XXX: No way to know the difference between the param being sent
	// as an empty string and the param not being sent. Create
	// a `Parameters` type that has an `Exists` method
	return req.raw.Form[name][0]
}

func (req *Request) EveryParam(name string) []string {
	return req.raw.Form[name]
}
