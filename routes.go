package mojo

import (
	"fmt"
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
	*Routes
	Name     string
	Methods  util.StringSlice
	Pattern  *regexp.Regexp
	Defaults Stash
	Handler  Handler
}

// Match is a set of route destinations for a given request
type Match struct {
	Stack    []*Route
	Position int
}

var stdPlaceholder *regexp.Regexp

func init() {
	stdPlaceholder = regexp.MustCompile("(^|[/<])([:])([a-zA-Z_]+)($|[/>])")
}

// Any creates a Route to handle any of the given methods for the given
// path. The path can contain placeholders which will populate values in
// the Stash.
//
// Standard placeholders begin with ":" and match all characters except
// for "/" and ".".
//
// Placeholders can be made optional by providing default values in
// a Stash as an additional argument. Slashes before optional
// placeholders also become optional.
//
// Placeholders can be enclosed in "<" and ">" to separate them from the
// surrounding text.
//
// If placeholders are not powerful enough, the path can contain named
// capture groups as a regular expression, like "(?P<name>\d+)" to match
// only digits.
func (rs *Routes) Any(methods []string, path string, opts ...interface{}) *Route {
	stash := optionalStash(opts)
	pathPattern := parsePattern(path, stash)
	r := &Route{
		Methods:  methods,
		Pattern:  regexp.MustCompile(pathPattern),
		Defaults: stash,
	}
	rs.routes = append(rs.routes, r)
	return r
}

// optionalStash finds and returns the first Stash object in the given
// array of options to a function. If no Stash object is found, returns
// an empty Stash object.
func optionalStash(opts []interface{}) Stash {
	// XXX: If we require a stash object, create a requireStash() func
	// XXX: Remove found stash from opts so we can assert that all
	// opts have been used
	for _, opt := range opts {
		if stash, ok := opt.(Stash); ok {
			return stash
		}
	}
	return Stash{}
}

// parsePattern parses the given path with placeholders and returns
// a string suitable for a regexp. This string does not contain
// start/end anchors, so that different route types can choose different
// anchors.
func parsePattern(path string, stash Stash) string {
	// Standard placeholders
	// XXX: Add relaxed and wildcard placeholders
	// XXX: Add restricted placeholders
	pathPattern := ""
	lastIndex := 0
	for _, v := range stdPlaceholder.FindAllStringSubmatchIndex(path, -1) {
		// v is an array of pairs of ints. Each pair is start and end
		// indexes of the match in the string. The first pair is the
		// entire match, and other pairs are the corresponding
		// submatches
		gap := path[lastIndex:v[0]]
		lastIndex = v[1]

		start := path[v[2]:v[3]]
		if start != "/" {
			start = ""
		}

		//placeType := path[v[4]:v[5]]
		placeName := path[v[6]:v[7]]
		//end := path[v[8]:v[9]] // unused

		matchType := "+" // required
		if _, ok := stash[placeName]; ok {
			matchType = "*" // optional
			if start == "/" {
				start += "?"
			}
		}

		pathPattern += gap + fmt.Sprintf("%s(?P<%s>[^/.]%s)", start, placeName, matchType)
	}

	// If we never matched, there were no placeholders
	if pathPattern == "" {
		return path
	}

	return pathPattern
}

func (rs *Routes) Get(pattern string, opts ...interface{}) *Route {
	return rs.Any([]string{"GET"}, pattern, opts...)
}
func (rs *Routes) Post(pattern string, opts ...interface{}) *Route {
	return rs.Any([]string{"POST"}, pattern, opts...)
}
func (rs *Routes) Put(pattern string, opts ...interface{}) *Route {
	return rs.Any([]string{"PUT"}, pattern, opts...)
}
func (rs *Routes) Patch(pattern string, opts ...interface{}) *Route {
	return rs.Any([]string{"PATCH"}, pattern, opts...)
}
func (rs *Routes) Delete(pattern string, opts ...interface{}) *Route {
	return rs.Any([]string{"DELETE"}, pattern, opts...)
}

// Under creates an intermediate destination. Under routes can have
// further destinations nested inside them. The handler given to Under
// must return a boolean to determine whether to continue dispatch.
func (rs *Routes) Under(pattern string, handler func(*Context) bool, opts ...interface{}) *Route {
	stash := optionalStash(opts)
	pathPattern := parsePattern(pattern, stash)
	r := &Route{
		Methods:  []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		Pattern:  regexp.MustCompile(pathPattern),
		Defaults: stash,
		Routes:   &Routes{},
		Handler: func(c *Context) {
			c.continueDispatch = handler(c)
		},
	}
	rs.routes = append(rs.routes, r)
	return r
}

// findRoute finds a matching route for the given method and path,
// returning the route and any placeholder values to be put into the
// stash
func (rs *Routes) findRoute(method string, path string) (*Route, *Stash) {
	for _, r := range rs.routes {
		// Check method
		if !r.Methods.Has(method) {
			continue
		}

		// Check path
		// BUG: If no children, must match completely (not prefix)
		regexpMatch := r.Pattern.FindStringSubmatch(path)
		if regexpMatch == nil {
			continue
		}

		// Add placeholder values to stash
		stash := Stash{}
		stash.Merge(r.Defaults)
		stashNames := r.Pattern.SubexpNames()
		for i, value := range regexpMatch {
			if i == 0 || value == "" {
				continue
			}
			stash[stashNames[i]] = value
		}
		return r, &stash
	}
	return nil, nil
}

// Match tries to find the route(s) for the given Context. If found, the
// context's Match value will be set to an array of Route objects to
// call, in order. See: Dispatch
func (rs *Routes) Match(c *Context) {
	method := c.Req.Method
	path := c.Stash["path"].(string)
	// XXX: Replace with Log
	//fmt.Printf("[debug] %s %s\n", method, path)

	r, stash := rs.findRoute(method, path)
	if r == nil {
		return
	}

	// Matched!
	if c.Match == nil {
		c.Match = &Match{}
	}
	c.Match.Append(r)

	// XXX: If there are child routes of this route (Under), try to
	// append to the match stack
	for r.Routes != nil && len(r.Routes.routes) > 0 {
		// Remove the last match off the front
		childPath := r.Pattern.ReplaceAllString(path, "")
		// Try to match again
		cr, cstash := r.findRoute(method, childPath)

		if cr == nil {
			// If we had a child to match but didn't, matching fails
			c.Match = nil
			return
		}
		c.Match.Append(cr)
		stash.Merge(*cstash)

		// Keep trying
		r = cr
		path = childPath
	}

	// Matching succeeded!
	c.Stash.Merge(*stash)
}

// Dispatch takes the given context, finds matching route(s), and calls
// their handlers.
func (rs *Routes) Dispatch(c *Context) {
	rs.Match(c)
	if c.Match == nil {
		// XXX: Replace with Log
		//fmt.Printf("[debug] 404 Not Found\n")
		c.Res.Code = 404
		c.Res.Status = "Not Found"
		return
	}
	// Call handler in matched Route objects
	// XXX: Does this handle async correctly?
	for _, r := range c.Match.Stack {
		// XXX: Create mojo.Log w/ Error, Warning, Info, Debug, Trace
		if r.Handler != nil {
			r.Handler(c)
		}
		if !c.continueDispatch {
			break
		}
	}
}

func (r *Route) To(handler Handler) *Route {
	r.Handler = handler
	return r
}

func (m *Match) Append(r *Route) {
	m.Stack = append(m.Stack, r)
}
