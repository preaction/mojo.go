package mojo

// Controller is the central object for request handling
type Controller struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
	Match *Match
}
