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
	var stash Stash
	for _, opt := range opts {
		switch v := opt.(type) {
		case Stash:
			stash = v
		default:
			panic(fmt.Sprintf("Invalid argument %T", v))
		}
	}

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

	r := &Route{
		Methods:  methods,
		Pattern:  regexp.MustCompile(pathPattern),
		Defaults: stash,
	}
	rs.routes = append(rs.routes, r)
	return r
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

func (rs *Routes) Match(c *Context) {
	method := c.Req.Method
	path := c.Stash["path"].(string)
	// XXX: Replace with Log
	//fmt.Printf("[debug] %s %s\n", method, path)

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
		// Add Route defaults to stash
		c.Stash.Merge(r.Defaults)
		// Add placeholder values to stash
		stashNames := r.Pattern.SubexpNames()
		for i, value := range regexpMatch {
			if i == 0 || value == "" {
				continue
			}
			c.Stash[stashNames[i]] = value
		}

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
	if c.Match == nil {
		// XXX: Replace with Log
		//fmt.Printf("[debug] 404 Not Found\n")
		c.Res.Code = 404
		c.Res.Message = "Not Found"
		return
	}
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

func (m *Match) Append(r *Route) {
	m.Stack = append(m.Stack, r)
}
