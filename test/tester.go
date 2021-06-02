package test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
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

// GetOk tries a GET request to the given path. This test passes if the
// request is completed without panicking.
func (t *Tester) GetOk(path string, name ...string) *Tester {
	t.T.Helper()
	fillName(&name, fmt.Sprintf("GET %d", path))

	// XXX: Create mojo.URL wrapper
	// Everything needs a damned wrapper to be even remotely usable...
	// XXX: Create a server to integration test (and in case we want to
	// turn Application and Server into interfaces in the future)
	url, _ := url.ParseRequestURI(path)
	c := mojo.Context{
		Req: &mojo.Request{Method: "GET", URL: url},
		Res: &mojo.Response{Writer: httptest.NewRecorder()},
	}
	t.App.BuildContext(&c)
	t.Context = &c

	defer func() {
		if r := recover(); r != nil {
			t.T.Errorf("Failed test '%s': Panic in route handler: %v", name[0], r)
			t.Success = false
		}
	}()
	t.App.Handler(&c)
	t.Success = true
	return t
}

// StatusIs tests the status code of the current request.
func (t *Tester) StatusIs(code int, name ...string) *Tester {
	t.T.Helper()
	fillName(&name, fmt.Sprintf("Status is %d", code))
	if !t.hasRes(name) {
		return t
	}

	if t.Context.Res.Code != code {
		t.T.Errorf("Failed test '%s': Status %d != %d", name[0], t.Context.Res.Code, code)
		t.Success = false
	} else {
		t.Success = true
	}
	return t
}

// hasRes returns true if there is a current request to test
func (t *Tester) hasRes(name []string) bool {
	if t.Context == nil || t.Context.Res == nil {
		t.T.Errorf("Failed test '%s': Status is nil", name[0])
		t.Success = false
		return false
	}
	return true
}

// fillName fills in a default name for a test if needed
func fillName(name *[]string, defaultName string) {
	if name == nil {
		*name = []string{defaultName}
	}
}
