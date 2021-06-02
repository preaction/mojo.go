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

	// XXX: Create mojo.URL wrapper
	// Everything needs a damned wrapper to be even remotely usable...
	url, err := url.ParseRequestURI("/foo")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	// XXX: Create helper to build controller from transaction
	// XXX: Create helper to build transaction from method/url
	c := mojo.Context{
		Req: &mojo.Request{Method: "GET", URL: url},
		Stash: map[string]interface{}{
			"path": url.Path,
		},
	}

	r.Dispatch(&c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
}
