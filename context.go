package mojo

// Stash is a place to store arbitrary data during a request.
type Stash map[string]interface{}

// Merge adds the keys from the given stash into the current stash. This
// is a shallow merge (for now).
func (s *Stash) Merge(src Stash) {
	for k, v := range src {
		(*s)[k] = v
	}
}

// Context is the central object for request handling
type Context struct {
	Req   *Request
	Res   *Response
	Stash Stash
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
