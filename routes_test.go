package mojo_test

import (
	"net/url"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestRoutesMatch(t *testing.T) {
	r := mojo.Routes{}
	r.Any([]string{"GET"}, "/foo/bar")
	r.Get("/foo")

	// XXX: Create mojo.URL wrapper
	// Everything needs a damned wrapper to be even remotely usable...
	url, err := url.ParseRequestURI("/foo")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	c := mojo.Controller{
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
