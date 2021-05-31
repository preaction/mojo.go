package mojo

type Controller struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
}
