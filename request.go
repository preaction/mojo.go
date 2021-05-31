package mojo

import "net/http"

type Request struct {
	Message
	raw *http.Request
}

func (req *Request) Method() string {
	return req.raw.Method
}

func (req *Request) Param(name string) string {
	// XXX: No way to know the difference between the param being sent
	// as an empty string and the param not being sent. Create
	// a `Parameters` type that has an `Exists` method
	return req.raw.Form[name][0]
}

func (req *Request) EveryParam(name string) []string {
	return req.raw.Form[name]
}
