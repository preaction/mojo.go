package mojo

// Context is the central object for request handling
type Context struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
	Match *Match
}

// Param returns the given parameter. Stash values take precedence over
// form values, which take precedence over query parameters. Returns the
// empty string if the parameter is not found. Panics if the stash value
// is not a string.
func (c *Context) Param(name string) string {
	if value, ok := c.Stash[name]; ok {
		return value.(string)
	}
	return c.Req.Param(name)
}
