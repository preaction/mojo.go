package test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/preaction/mojo.go"
)

type Tester struct {
	*testing.T
	App     *mojo.Application
	Success bool
	Context *mojo.Context
}

func (t *Tester) GetOk(path string, name ...string) *Tester {
	t.T.Helper()
	if name == nil {
		name = []string{fmt.Sprintf("GET %d", path)}
	}

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

func (t *Tester) StatusIs(code int, name ...string) *Tester {
	t.T.Helper()
	if name == nil {
		name = []string{fmt.Sprintf("Status is %d", code)}
	}
	if t.Context == nil || t.Context.Res == nil {
		t.T.Errorf("Failed test '%s': Status is nil", name[0])
		t.Success = false
	} else if t.Context.Res.Code != code {
		t.T.Errorf("Failed test '%s': Status %d != %d", name[0], t.Context.Res.Code, code)
		t.Success = false
	} else {
		t.Success = true
	}
	return t
}
