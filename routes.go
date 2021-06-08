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
	Defaults map[string]interface{}
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
// Placeholders can be enclosed in "<" and ">" to separate them from the
// surrounding text.
//
// If placeholders are not powerful enough, the path can contain named
// capture groups as a regular expression, like "(?P<name>\d+)" to match
// only digits.
func (rs *Routes) Any(methods []string, path string) *Route {
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
		//placeType := path[v[4]:v[5]]
		placeName := path[v[6]:v[7]]
		end := path[v[8]:v[9]]

		pathPattern += gap + fmt.Sprintf("%s(?P<%s>[^/.]+)%s", start, placeName, end)
	}

	r := &Route{Methods: methods, Pattern: regexp.MustCompile(pathPattern)}
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
		// XXX: Add Route defaults to stash
		// Add placeholder values to stash
		stashNames := r.Pattern.SubexpNames()
		for i, value := range regexpMatch {
			if i == 0 {
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

func (r *Route) Stash(defaults map[string]interface{}) {
	for key, val := range defaults {
		r.Defaults[key] = val
	}
}

func (m *Match) Append(r *Route) {
	m.Stack = append(m.Stack, r)
}
