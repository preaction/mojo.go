package mojo_test

import (
	"net/url"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestRoutesMatch(t *testing.T) {
	r := mojo.Routes{}
	r.Get("/foo")

	// XXX: Create mojo.URL wrapper
	// Everything needs a damned wrapper to be even remotely usable...
	url, err := url.ParseRequestURI("/foo")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	c := mojo.Context{
		Req: &mojo.Request{Method: "GET", URL: url},
		Stash: map[string]interface{}{
			"path": url.Path,
		},
	}

	r.Dispatch(&c)

	if c.Match == nil {
		t.Errorf("No route found for request: %v", c.Req)
	}
}

func TestRoutesHandler(t *testing.T) {
	handlerCalled := false
	handler := func(c *mojo.Context) {
		handlerCalled = true
	}

	r := mojo.Routes{}
	r.Get("/foo").To(handler)

	// XXX: Create helper to build controller from transaction
	c := mojo.Context{
		Req: mojo.NewRequest("GET", "/foo"),
		Stash: map[string]interface{}{
			"path": "/foo",
		},
	}

	r.Dispatch(&c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
}

func TestRoutesPlaceholder(t *testing.T) {
	name := ""
	handlerCalled := false
	handler := func(c *mojo.Context) {
		name = c.Stash["name"].(string)
		handlerCalled = true
	}

	r := mojo.Routes{}
	r.Get("/foo/:name").To(handler)

	// XXX: Create helper to build controller from transaction
	c := mojo.Context{
		Req: mojo.NewRequest("GET", "/foo/morbo"),
		Stash: map[string]interface{}{
			"path": "/foo/morbo",
		},
	}

	r.Dispatch(&c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
	if name != "morbo" {
		t.Errorf("Route placeholder not added to stash")
	}
}
