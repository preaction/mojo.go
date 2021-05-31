package mojo

import (
	"regexp"

	"github.com/preaction/mojo.go/util"
)

// Routes stores the application routes and handles matching routes to
// incoming requests
type Routes struct {
	routes []Route
}

// Route is a single endpoint
type Route struct {
	Name     string
	Methods  util.StringSlice
	Pattern  *regexp.Regexp
	Defaults map[string]interface{}
}

// Match is a set of route destinations for a given request
type Match struct {
	Stack    []*Route
	Position int
}

func (rs *Routes) Any(methods []string, pattern string) *Route {
	// XXX: Transform pattern to regexp with named placeholders
	r := Route{Methods: methods, Pattern: regexp.MustCompile(pattern)}
	rs.routes = append(rs.routes, r)
	return &r
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

func (rs *Routes) Match(c *Controller) {
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

		// Matched!
		if c.Match == nil {
			c.Match = &Match{}
		}
		c.Match.Append(&r)

		// XXX: If there are child routes of this route (Under), try to
		// append to the match stack
	}
}

func (rs *Routes) Dispatch(c *Controller) {
	rs.Match(c)
	// XXX: Call handler in matched Route objects
}

func (r *Route) To(defaults map[string]interface{}) *Route {
	for key, val := range defaults {
		r.Defaults[key] = val
	}
	return r
}

func (m *Match) Append(r *Route) {
	m.Stack = append(m.Stack, r)
}
