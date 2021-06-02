package mojo

import (
	"regexp"

	"github.com/preaction/mojo.go/util"
)

// Routes stores the application routes and handles matching routes to
// incoming requests
type Routes struct {
	routes []*Route
}

// Handler handles an incoming request. Handlers can modify the stash or
// response for subsequent handlers or render a response body.
type Handler func(*Context)

// Route is a single endpoint
type Route struct {
	Name     string
	Methods  util.StringSlice
	Pattern  *regexp.Regexp
	Defaults map[string]interface{}
	Handler  Handler
}

// Match is a set of route destinations for a given request
type Match struct {
	Stack    []*Route
	Position int
}

func (rs *Routes) Any(methods []string, pattern string) *Route {
	// XXX: Transform pattern to regexp with named placeholders
	r := &Route{Methods: methods, Pattern: regexp.MustCompile(pattern)}
	rs.routes = append(rs.routes, r)
	return r
}

func (rs *Routes) Get(pattern string) *Route {
	return rs.Any([]string{"GET"}, pattern)
}
func (rs *Routes) Post(pattern string) *Route {
	return rs.Any([]string{"POST"}, pattern)
}
func (rs *Routes) Put(pattern string) *Route {
	return rs.Any([]string{"PUT"}, pattern)
}
func (rs *Routes) Patch(pattern string) *Route {
	return rs.Any([]string{"PATCH"}, pattern)
}
func (rs *Routes) Delete(pattern string) *Route {
	return rs.Any([]string{"DELETE"}, pattern)
}

func (rs *Routes) Match(c *Context) {
	method := c.Req.Method
	path := c.Stash["path"].(string)

	for _, r := range rs.routes {
		// Check method
		if !r.Methods.Has(method) {
			continue
		}

		// Check path
		regexpMatch := r.Pattern.FindStringSubmatch(path)
		if regexpMatch == nil {
			continue
		}
		// XXX: Add Route defaults to stash
		// XXX: Add placeholder values to stash

		// Matched!
		if c.Match == nil {
			c.Match = &Match{}
		}
		c.Match.Append(r)

		// XXX: If there are child routes of this route (Under), try to
		// append to the match stack
	}
}

func (rs *Routes) Dispatch(c *Context) {
	rs.Match(c)
	// Call handler in matched Route objects
	// XXX: Does this handle async correctly?
	for _, r := range c.Match.Stack {
		// XXX: Create mojo.Log w/ Error, Warning, Info, Debug, Trace
		if r.Handler != nil {
			r.Handler(c)
		}
	}
}

func (r *Route) To(handler Handler) *Route {
	r.Handler = handler
	return r
}

func (r *Route) Stash(defaults map[string]interface{}) {
	for key, val := range defaults {
		r.Defaults[key] = val
	}
}

func (m *Match) Append(r *Route) {
	m.Stack = append(m.Stack, r)
}
