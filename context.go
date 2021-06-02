package mojo

// Context is the central object for request handling
type Context struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
	Match *Match
}
