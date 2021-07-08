package test

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/preaction/mojo.go"
)

// Tester is a helper for testing Mojo Applications. Tester wraps
// a testing.T object and keeps the state of the current request being
// tested.
type Tester struct {
	*testing.T
	App     *mojo.Application
	Success bool
	Context *mojo.Context
}

// NewTester creates a new tester for the given application
func NewTester(t *testing.T, app *mojo.Application) *Tester {
	return &Tester{T: t, App: app}
}

// NewContext returns a new context with sensible defaults for testing.
// Any opts provided will override the defaults for testing.
func NewContext(t *testing.T, opts ...interface{}) *mojo.Context {
	c := mojo.Context{
		Req:   mojo.NewRequest("GET", "/"),
		Res:   mojo.NewResponse(httptest.NewRecorder()),
		Stash: map[string]interface{}{},
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case *mojo.Request:
			c.Req = v
		case *mojo.Response:
			c.Res = v
		case mojo.Stash:
			c.Stash.Merge(v)
		default:
			t.Fatalf("Unknown option type %T\n", v)
		}
	}

	// Set default stash from request
	c.Stash["path"] = c.Req.URL.Path

	return &c
}

// BuildHTTPRequest builds a new http.Request object from the given raw HTTP
// request
func BuildHTTPRequest(t *testing.T, raw string) *http.Request {
	t.Helper()
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatalf("Could not read request: %v", err)
	}
	return req
}

// ReadHTTPResponse returns the response content from the given test context
// (see NewContext). It is assumed the internal http.Response is an
// httptest.ResponseRecorder.
func ReadHTTPResponse(t *testing.T, c *mojo.Context) (*httptest.ResponseRecorder, []byte) {
	t.Helper()
	res := c.Res.Writer.(*httptest.ResponseRecorder)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %v", err)
	}
	return res, body
}

// GetOk tries a GET request to the given path. This test passes if the
// request is completed without panicking.
func (t *Tester) GetOk(path string, name ...string) *Tester {
	t.T.Helper()
	fillName(&name, fmt.Sprintf("GET %d", path))

	// XXX: Create a server to integration test (and in case we want to
	// turn Application and Server into interfaces in the future)
	req := mojo.NewRequest("GET", path)
	res := mojo.NewResponse(httptest.NewRecorder())
	c := t.App.BuildContext(req, res)
	t.Context = c

	defer func() {
		if r := recover(); r != nil {
			t.errorf(name, "Panic in route handler: %v", r)
		}
	}()
	t.Success = true
	t.App.Handler(c)
	return t
}

// errorf prints the formatted error and updates the Success flag
func (t *Tester) errorf(name []string, text string, args ...interface{}) {
	t.T.Helper()
	t.T.Errorf("Failed test '%s': %s", name[0], fmt.Sprintf(text, args...))
	t.Success = false
}

// StatusIs tests the status code of the current request.
func (t *Tester) StatusIs(code int, name ...string) *Tester {
	t.T.Helper()
	fillName(&name, fmt.Sprintf("Status is %d", code))
	if !t.hasRes(name) {
		return t
	}

	if t.Context.Res.Code != code {
		t.errorf(name, "Status %d != %d", t.Context.Res.Code, code)
		return t
	}

	t.Success = true
	return t
}

// TextIs tests the response body of the current request equals the
// given text
func (t *Tester) TextIs(text string, name ...string) *Tester {
	t.T.Helper()
	fillName(&name, "Response text")
	if !t.hasRes(name) {
		return t
	}

	if t.Context.Res.Content.String() != text {
		t.errorf(name, "Text is not equal:\n\tExpect: %s\n\tGot: %s", text, t.Context.Res.Content.String())
		return t
	}
	t.Success = true
	return t
}

// hasRes returns true if there is a current request to test, updating
// the Success flag if not
func (t *Tester) hasRes(name []string) bool {
	if t.Context == nil || t.Context.Res == nil {
		t.errorf(name, "Status is nil")
		return false
	}
	return true
}

// fillName fills in a default name for a test if needed
func fillName(name *[]string, defaultName string) {
	if name == nil || len(*name) == 0 {
		*name = []string{defaultName}
	}
}
